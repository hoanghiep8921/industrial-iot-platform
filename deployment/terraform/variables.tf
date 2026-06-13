# ============================================================
# Industrial IoT Platform - Terraform Variables
# ============================================================

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "industrial-iot"
}

variable "environment" {
  description = "Deployment environment (dev, staging, production)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "production"], var.environment)
    error_message = "Environment must be one of: dev, staging, production."
  }
}

variable "cloud_provider" {
  description = "Cloud provider (aws, gcp, azure)"
  type        = string
  default     = "aws"
}

variable "regions" {
  description = "Primary region for cloud deployment"
  type        = string
  default     = "ap-southeast-1" # Singapore (gần Việt Nam)
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.29"
}

variable "node_groups" {
  description = "EKS node group configuration"
  type = map(object({
    instance_type = string
    min_size      = number
    max_size      = number
    desired_size  = number
    labels        = map(string)
  }))
  default = {
    general = {
      instance_type = "t3.xlarge"
      min_size      = 2
      max_size      = 10
      desired_size  = 3
      labels = {
        role = "general"
      }
    }
    database = {
      instance_type = "r6g.xlarge"
      min_size      = 1
      max_size      = 3
      desired_size  = 2
      labels = {
        role = "database"
      }
    }
    mqtt = {
      instance_type = "c6g.xlarge"
      min_size      = 2
      max_size      = 6
      desired_size  = 3
      labels = {
        role = "mqtt"
      }
    }
  }
}

variable "mqtt_lb_ports" {
  description = "MQTT load balancer ports"
  type = map(number)
  default = {
    mqtt      = 1883
    mqtts     = 8883
    websocket = 8083
  }
}

variable "database_config" {
  description = "RDS database configuration"
  type = object({
    engine              = string
    engine_version      = string
    instance_class      = string
    allocated_storage   = number
    max_allocated_storage = number
    multi_az            = bool
    backup_retention_days = number
  })
  default = {
    engine                 = "postgres"
    engine_version         = "16"
    instance_class         = "db.r6g.xlarge"
    allocated_storage      = 100
    max_allocated_storage  = 1000
    multi_az               = true
    backup_retention_days  = 30
  }
}

variable "edge_gateway_count" {
  description = "Number of edge gateways per region"
  type = object({
    hanoi     = number
    hochiminh = number
    danang    = number
  })
  default = {
    hanoi     = 2
    hochiminh = 2
    danang    = 2
  }
}

variable "tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    Project     = "industrial-iot"
    ManagedBy   = "terraform"
    Owner       = "engineering"
  }
}
