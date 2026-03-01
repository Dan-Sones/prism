output "app_ip" {
  value = google_compute_instance.app.network_interface[0].access_config[0].nat_ip
}

output "gatling_ip" {
  value = google_compute_instance.gatling.network_interface[0].access_config[0].nat_ip
}
