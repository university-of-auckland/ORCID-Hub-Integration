# TODO! Need to figure out how to access Kong admin API
# TODO! Or perhaps the the consumer also has be created as part of this excersize...
# Kong entry point -> AWS API Gateway
# TODO: upstream_url should be taken from the AWS API gateway resource...
resource "null_resource" "kong-api-setup" {
	depends_on = ["null_resource.connector"]
  provisioner "local-exec" {
      command = <<EOT
        curl  -v > kong-api-setup.log \
					-X POST https://api.dev.auckland.ac.nz/admin/apis/ \
					-H 'Content-Type: application/json' \
					-H 'Accept: application/json' \
					-H 'apikey: ${var.apikey}' \
					-d '{ \
						"http_if_terminated": false, "https_only": false, "name": "orcidhub-trigger-integration-v1", \
						"preserve_host": false, \
						"retries": 5, \
						"strip_uri": true, \
						"upstream_connect_timeout": 60000, \
						"upstream_read_timeout": 60000, \
						"upstream_send_timeout": 60000, \
						"upstream_url": " https://415mdw939a.execute-api.ap-southeast-2.amazonaws.com/dev/v1/enqueue", \
						"methods": ["POST"], \
						"uris": ["/orcidhub-trigger/integrations/v1"] \
					}'
EOT
  }
}

resource "null_resource" "kong-api-cleanup" {
  provisioner "local-exec" {
		  when = "destroy"
      command = <<EOT
        curl  -v > kong-api-cleanup.log \
					-X DELETE https://api.dev.auckland.ac.nz/admin/apis/orcidhub-trigger-integration-v1 \
					-H 'Accept: application/json' \
					-H 'apikey: ${var.apikey}'
EOT
  }
}

# Kafka connector: Kafka -> Webhook Kong entry point
resource "null_resource" "connector" {
  provisioner "local-exec" {
      command = <<EOT
        curl  -v > connection-orcidhub-setup.log \
					-X POST https://api.dev.auckland.ac.nz/service/kafka-connect/v2/connectors/ \
					-H 'Content-Type: application/json' \
					-H 'Accept: application/json' \
					-H 'apikey: ${var.apikey}' \
					-d '{ \
						"name": "connection-orcidhub", \
						"config": { \
								"connector.class":"nz.ac.auckland.kafka.http.sink.HttpSinkConnector", \
								"tasks.max":"1", \
								"value.converter":"org.apache.kafka.connect.storage.StringConverter", \
								"key.converter":"org.apache.kafka.connect.storage.StringConverter", \
								"header.converter":"org.apache.kafka.connect.storage.StringConverter", \
								"topics":"nz-ac-auckland-employment", \
								"callback.request.url":"https://api.dev.auckland.ac.nz/service/orcidhub-trigger/integrations/v1", \
								"callback.request.method":"POST", \
								"callback.request.headers":"apikey:${var.apikey}|Content-Type:application/json", \
								"retry.backoff.sec":"600,86400,259200", \
								"exception.strategy":"PROGRESS_BACK_OFF_DROP_MESSAGE" \
						} \
					}'
EOT
  }
}

resource "null_resource" "connector-cleanup" {
  provisioner "local-exec" {
		  when = "destroy"
      command = <<EOT
        curl  -v > connection-orcidhub-cleanup.log \
					-X DELETE https://api.dev.auckland.ac.nz/service/kafka-connect/v2/connectors/connection-orcidhub \
					-H 'Accept: application/json' \
					-H 'apikey: ${var.apikey}'
EOT
  }
}
