mockgen -source ../api/http_client/my_http_client.go -destination ../testing/mocks/api/http_client/mock_my_http_client.go 
mockgen -source ../api/keel/keel_webhook_api.go -destination ../testing/mocks/api/keel/mock_keel_webhook_api.go
mockgen -source ../api/kubeclarity/kubeclarity_runtime_api.go -destination ../testing/mocks/api/kubeclarity/mock_kubeclarity_runtime_api.go
mockgen -source ../api/kubernetes/kubernetes_api.go -destination ../testing/mocks/api/kubernetes/mock_kubernetes_api.go
mockgen -source ../service/kubernetes/kubernetes_resource_service.go -destination ../testing/mocks/service/kubernetes/mock_kubernetes_resource_service.go
mockgen -source ../service/resource/kubeclarity_resource_service.go -destination ../testing/mocks/service/resource/mock_kubeclarity_resource_service.go
mockgen -source ../service/scan/kubeclarity_scan_service.go -destination ../testing/mocks/service/scan/mock_kubeclarity_scan_service.go
mockgen -source ../service/update/keel_update_service.go -destination ../testing/mocks/service/update/mock_keel_update_service.go
mockgen -source ../service/notification/email_notification_service.go -destination ../testing/mocks/service/notification/mock_email_notification_service.go