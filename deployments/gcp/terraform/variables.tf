# Terraform variables for Portfolio Backend infrastructure

variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region for resources"
  type        = string
  default     = "us-central1"
}

variable "db_name" {
  description = "Name of the portfolio database"
  type        = string
  default     = "portfolio_db"
}

variable "db_username" {
  description = "Username for the portfolio database"
  type        = string
  default     = "portfolio_user"
}

variable "db_password" {
  description = "Password for the portfolio database user"
  type        = string
  sensitive   = true
}

variable "db_tier" {
  description = "Machine type for Cloud SQL instance"
  type        = string
  default     = "db-f1-micro"
  
  validation {
    condition = contains([
      "db-f1-micro",
      "db-g1-small",
      "db-n1-standard-1",
      "db-n1-standard-2",
      "db-n1-standard-4"
    ], var.db_tier)
    error_message = "Database tier must be a valid Cloud SQL machine type."
  }
}

variable "db_availability_type" {
  description = "Availability type for Cloud SQL instance"
  type        = string
  default     = "ZONAL"
  
  validation {
    condition     = contains(["ZONAL", "REGIONAL"], var.db_availability_type)
    error_message = "Availability type must be either ZONAL or REGIONAL."
  }
}

variable "db_disk_size" {
  description = "Disk size in GB for Cloud SQL instance"
  type        = number
  default     = 10
  
  validation {
    condition     = var.db_disk_size >= 10 && var.db_disk_size <= 30720
    error_message = "Disk size must be between 10 and 30720 GB."
  }
}

variable "deletion_protection" {
  description = "Enable deletion protection for Cloud SQL instance"
  type        = bool
  default     = true
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "prod"
  
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "cors_allowed_origins" {
  description = "List of allowed CORS origins"
  type        = list(string)
  default     = ["https://your-frontend-domain.com"]
}