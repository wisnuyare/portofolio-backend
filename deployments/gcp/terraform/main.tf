# Terraform configuration for Portfolio Backend GCP infrastructure

terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "cloudbuild.googleapis.com",
    "secretmanager.googleapis.com",
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
  ])

  project = var.project_id
  service = each.value

  disable_dependent_services = true
}

# Create VPC network for private services
resource "google_compute_network" "vpc" {
  name                    = "portfolio-vpc"
  auto_create_subnetworks = false

  depends_on = [google_project_service.required_apis]
}

# Create subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "portfolio-subnet"
  ip_cidr_range = "10.0.0.0/24"
  region        = var.region
  network       = google_compute_network.vpc.id
}

# Reserve IP range for private services
resource "google_compute_global_address" "private_ip_range" {
  name          = "portfolio-private-ip-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

# Create private connection
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_range.name]
}

# Create Cloud SQL instance
resource "google_sql_database_instance" "portfolio_db" {
  name             = "portfolio-db"
  database_version = "MYSQL_8_0"
  region           = var.region

  settings {
    tier              = var.db_tier
    availability_type = var.db_availability_type
    disk_size         = var.db_disk_size
    disk_type         = "PD_SSD"

    # Backup configuration
    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true
      backup_retention_settings {
        retained_backups = 7
      }
    }

    # IP configuration for private access
    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = google_compute_network.vpc.id
      enable_private_path_for_google_cloud_services = true
    }

    # Maintenance window
    maintenance_window {
      day          = 7  # Sunday
      hour         = 4  # 4 AM
      update_track = "stable"
    }

    # Database flags
    database_flags {
      name  = "innodb_buffer_pool_size"
      value = "134217728"  # 128MB
    }
  }

  deletion_protection = var.deletion_protection

  depends_on = [
    google_service_networking_connection.private_vpc_connection,
    google_project_service.required_apis
  ]
}

# Create database
resource "google_sql_database" "portfolio" {
  name     = var.db_name
  instance = google_sql_database_instance.portfolio_db.name
}

# Create database user
resource "google_sql_user" "portfolio_user" {
  name     = var.db_username
  instance = google_sql_database_instance.portfolio_db.name
  password = var.db_password
}

# Create service account for Cloud Run
resource "google_service_account" "portfolio_backend_sa" {
  account_id   = "portfolio-backend-sa"
  display_name = "Portfolio Backend Service Account"
  description  = "Service account for Portfolio Backend Cloud Run service"
}

# Grant Cloud SQL Client role to service account
resource "google_project_iam_member" "cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.portfolio_backend_sa.email}"
}

# Grant Secret Manager Secret Accessor role
resource "google_project_iam_member" "secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.portfolio_backend_sa.email}"
}

# Create database credentials secrets
resource "google_secret_manager_secret" "db_username" {
  secret_id = "db-username"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required_apis]
}

resource "google_secret_manager_secret_version" "db_username" {
  secret      = google_secret_manager_secret.db_username.id
  secret_data = var.db_username
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required_apis]
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.db_password
}

# Cloud Run service (basic configuration - detailed config in service.yaml)
resource "google_cloud_run_service" "portfolio_backend" {
  name     = "portfolio-backend"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.portfolio_backend_sa.email
      
      containers {
        image = "gcr.io/${var.project_id}/portfolio-backend:latest"
        
        ports {
          container_port = 8080
        }

        env {
          name  = "DB_HOST"
          value = "127.0.0.1"
        }
        
        env {
          name  = "DB_NAME"
          value = var.db_name
        }

        resources {
          limits = {
            cpu    = "1000m"
            memory = "512Mi"
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"                    = "100"
        "run.googleapis.com/cloudsql-instances"               = google_sql_database_instance.portfolio_db.connection_name
        "run.googleapis.com/execution-environment"            = "gen2"
        "run.googleapis.com/cpu-throttling"                   = "false"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.required_apis,
    google_sql_database_instance.portfolio_db
  ]
}

# Allow unauthenticated invocations
resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.portfolio_backend.location
  project  = google_cloud_run_service.portfolio_backend.project
  service  = google_cloud_run_service.portfolio_backend.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}