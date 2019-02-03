# Go Arm Stackdriver Temps (gast)

The goal of this project is to submit temperatures of Arm based devices like Raspberry Pis and ODroids to Google Cloud Platform's Stackdriver custom monitoring for alerting and reporting.

- Written in Golang
- Requires a credentials json file (see <a href="https://cloud.google.com/docs/authentication/production#auth-cloud-implicit-go">https://cloud.google.com/docs/authentication/production#auth-cloud-implicit-go</a>)
- Must have permission to write to stackdriver monitoring.
- Works on Arm and Linux
- Low memory and processor usage
- See live <a href="https://public.google.stackdriver.com/public/chart/14868659815927352142?drawMode=color&showLegend=true&theme=dark">https://public.google.stackdriver.com/public/chart/14868659815927352142?drawMode=color&showLegend=true&theme=dark</a>
