#!/bin/bash

BASE_URL="http://localhost:8081/api/events-catalog"

curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Shipped",
    "eventKey": "order_shipped",
    "description": "",
    "fields": [
      {
        "name": "Final Order Total",
        "fieldKey": "final_order_total",
        "dataType": "string"
      },
      {
        "name": "Order Total Without Discounts",
        "fieldKey": "order_total_without_discounts",
        "dataType": "string"
      },
      {
        "name": "Postage Total",
        "fieldKey": "postage_total",
        "dataType": "string"
      }
    ]
  }' | jq .
