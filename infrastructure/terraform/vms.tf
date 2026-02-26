resource "google_compute_instance" "app" {
  name         = "prism-app"
  machine_type = "e2-standard-4"
  tags         = ["prism"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2404-lts-amd64"
      size  = 50
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  metadata = {
    ssh-keys = "ubuntu:${file(var.ssh_public_key)}"
  }
}


resource "google_compute_instance" "gatling" {
  name         = "prism-gatling"
  machine_type = "e2-highcpu-4"
  tags         = ["prism"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2404-lts-amd64"
      size  = 20
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  metadata = {
    ssh-keys = "ubuntu:${file(var.ssh_public_key)}"
  }
}
