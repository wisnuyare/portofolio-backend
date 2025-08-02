#!/bin/bash

# Deploy script for Portfolio Backend to GCP
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID=${PROJECT_ID:-""}
REGION=${REGION:-"us-central1"}
SERVICE_NAME="portfolio-backend"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if gcloud is installed
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI is not installed. Please install it first."
        exit 1
    fi
    
    # Check if docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install it first."
        exit 1
    fi
    
    # Check if project ID is set
    if [ -z "$PROJECT_ID" ]; then
        PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
        if [ -z "$PROJECT_ID" ]; then
            log_error "PROJECT_ID is not set. Please set it as an environment variable or configure gcloud."
            exit 1
        fi
    fi
    
    log_info "Project ID: $PROJECT_ID"
    log_info "Region: $REGION"
}

authenticate() {
    log_info "Configuring authentication..."
    gcloud auth configure-docker
}

build_and_push() {
    log_info "Building and pushing Docker image..."
    
    # Build image
    docker build -f deployments/docker/Dockerfile -t gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest .
    
    # Push image
    docker push gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest
}

deploy_to_cloud_run() {
    log_info "Deploying to Cloud Run..."
    
    gcloud run deploy ${SERVICE_NAME} \
        --image gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest \
        --region ${REGION} \
        --platform managed \
        --allow-unauthenticated \
        --port 8080 \
        --memory 512Mi \
        --cpu 1 \
        --max-instances 100 \
        --set-env-vars LOG_LEVEL=info,LOG_FORMAT=json \
        --project ${PROJECT_ID}
}

run_migrations() {
    log_info "Running database migrations..."
    log_warn "Migration functionality needs to be implemented based on your migration strategy"
    # TODO: Implement migration logic
    # This could involve:
    # 1. Creating a migration job
    # 2. Running migrations via Cloud SQL proxy
    # 3. Using a separate migration container
}

get_service_url() {
    log_info "Getting service URL..."
    SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} \
        --region ${REGION} \
        --project ${PROJECT_ID} \
        --format 'value(status.url)')
    
    log_info "Service deployed successfully!"
    log_info "URL: ${SERVICE_URL}"
}

# Main deployment flow
main() {
    log_info "Starting deployment of Portfolio Backend..."
    
    check_prerequisites
    authenticate
    build_and_push
    deploy_to_cloud_run
    run_migrations
    get_service_url
    
    log_info "Deployment completed successfully!"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --project)
            PROJECT_ID="$2"
            shift 2
            ;;
        --region)
            REGION="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [--project PROJECT_ID] [--region REGION]"
            echo "  --project  GCP Project ID (default: current gcloud project)"
            echo "  --region   GCP Region (default: us-central1)"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Run main function
main