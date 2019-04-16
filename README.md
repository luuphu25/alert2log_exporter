# "# alert2log_exporter"

[x] Handle alert send to address webhook

[x] Store alert to elastic

[x] Capture prometheus request

[ ] Expose metrics to "/metrics"

[x] Send request to prometheus (time series in past)

[x] Edit data before store in elastic

[] Check error metrics not exists
    get metric from log 


# Nhận alerts từ alertmanager
- Cấu hình alertmanager gửi tới địa chỉ đang lắng nghe dưới dạng webhook-config (/webhook path)

# Bắt request prometheus gửi cho alertmanager
- Thêm địa chỉ của tool đang lắng nghe vào `alertmangers` của prometheus.yml ( as alertmanager) 
- Prometheus tạo alert bằng request POST cho alertmanager bằng api nên tool này handle ở path /api/v1/alerts. 

*Note: 
- Main flie : http_server.go


