resource "google_container_cluster" "cluster" {
  provider = "google-beta"
  project  = "${var.project}"

  name   = "${var.cluster_name}"
  region = "${var.region}"

  additional_zones = "${length(var.zones) > 0 ? var.zones : null}"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true

  initial_node_count = 1

  ip_allocation_policy {
    create_subnetwork = true
  }

  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }

    http_load_balancing {
      disabled = true
    }
  }

  workload_identity_config {
    identity_namespace = "${var.project}.svc.id.goog"
  }

  # Setting an empty username and password explicitly disables basic auth
  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

resource "google_container_node_pool" "system" {
  provider           = "google-beta"
  project            = "${var.project}"
  name               = "system"
  cluster            = "${google_container_cluster.cluster.name}"
  region             = "${var.region}"
  initial_node_count = 1

  autoscaling {
    min_node_count = 1
    max_node_count = 3
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type = "${var.system_node_type}"
    preemptible  = "${var.preemptible}"
    disk_size_gb = 50
    image_type   = "COS"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/servicecontrol",
    ]

    labels = {
      "node-role.stack.presslabs.org/presslabs-sys" = ""
    }

    taint {
      key    = "CriticalAddonsOnly"
      value  = "true"
      effect = "${var.system_node_taint_effect}"
    }

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }
  }
}

resource "google_container_node_pool" "database" {
  provider           = "google-beta"
  project            = "${var.project}"
  name               = "database"
  cluster            = "${google_container_cluster.cluster.name}"
  region             = "${var.region}"
  initial_node_count = 0

  autoscaling {
    min_node_count = 0
    max_node_count = 3
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type = "${var.database_node_type}"
    preemptible  = "${var.preemptible}"
    disk_size_gb = 50
    image_type   = "COS"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/servicecontrol",
    ]

    labels = {
      "node-role.stack.presslabs.org/database"  = ""
      "node-role.stack.presslabs.org/mysql"     = ""
      "node-role.stack.presslabs.org/memcached" = ""
    }

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }
  }
}

resource "google_container_node_pool" "wordpress" {
  provider           = "google-beta"
  project            = "${var.project}"
  name               = "wordpress"
  cluster            = "${google_container_cluster.cluster.name}"
  region             = "${var.region}"
  initial_node_count = 0

  autoscaling {
    min_node_count = 0
    max_node_count = 5
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type = "${var.wordpress_node_type}"
    preemptible  = "${var.preemptible}"
    disk_size_gb = 100
    image_type   = "COS"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/servicecontrol",
    ]

    labels = {
      "node-role.stack.presslabs.org/wordpress" = ""
    }

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }
  }
}

resource "google_container_node_pool" "wordpress_preemptible" {
  provider = "google-beta"
  project  = "${var.project}"

  # create preemptible-specific group only if we don't want the entire cluster to be preemptible
  # (eg. for dev)
  count = "${var.preemptible ? 0 : 1}"

  name               = "wordpress-preemptible"
  cluster            = "${google_container_cluster.cluster.name}"
  region             = "${var.region}"
  initial_node_count = 0

  autoscaling {
    min_node_count = 0
    max_node_count = 5
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type = "${var.wordpress_node_type}"
    preemptible  = true
    disk_size_gb = 100
    image_type   = "COS"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/servicecontrol",
    ]

    labels = {
      "node-role.stack.presslabs.org/wordpress" = ""
    }

    taint {
      key    = "cloud.google.com/gke-preemptible"
      value  = "true"
      effect = "NO_SCHEDULE"
    }

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }
  }
}
