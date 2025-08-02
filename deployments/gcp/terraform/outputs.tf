# Terraform outputs for Portfolio Backend infrastructure

output "cloud_run_url" {
  description = "URL of the deployed Cloud Run service"
  value       = google_cloud_run_service.portfolio_backend.status[0].url
}

output "cloud_run_service_name" {
  description = "Name of the Cloud Run service"
  value       = google_cloud_run_service.portfolio_backend.name
}

output "database_instance_name" {
  description = "Name of the Cloud SQL instance"
  value       = google_sql_database_instance.portfolio_db.name
}

output "database_connection_name" {
  description = "Connection name for Cloud SQL instance"
  value       = google_sql_database_instance.portfolio_db.connection_name
}

output "database_private_ip" {
  description = "Private IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.portfolio_db.private_ip_address
}

output "service_account_email" {
  description = "Email of the service account used by Cloud Run"
  value       = google_service_account.portfolio_backend_sa.email
}

output "vpc_network_name" {
  description = "Name of the VPC network"
  value       = google_compute_network.vpc.name
}

output "vpc_subnet_name" {
  description = "Name of the VPC subnet"
  value       = google_compute_subnetwork.subnet.name
}

output "db_secrets" {
  description = "Secret Manager secret names for database credentials"
  value = {
    username = google_secret_manager_secret.db_username.secret_id
    password = google_secret_manager_secret.db_password.secret_id
  }
  sensitive = true
}

output "project_id" {
  description = "GCP Project ID"
  value       = var.project_id
}

output "region" {
  description = "GCP Region"
  value       = var.region
}