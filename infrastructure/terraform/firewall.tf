resource "google_compute_firewall" "allow_ssh" {
  name    = "prism-allow-ssh"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["${var.allowed_ip}/32"]
  target_tags   = ["prism"]
}

resource "google_compute_firewall" "allow_http" {
  name    = "prism-allow-http"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["${var.allowed_ip}/32", "${google_compute_instance.gatling.network_interface[0].access_config[0].nat_ip}/32"]
  target_tags   = ["prism"]
}
