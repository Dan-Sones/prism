WITH order_shipped AS (
    INSERT INTO prism.event_types (name, event_key, version, description)
    VALUES ('Order Shipped', 'order_shipped', 1, 'Fired when an order is shipped')
    RETURNING id
)
INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type)
VALUES ((SELECT id FROM order_shipped), 'Final Order Total', 'final_order_total', 'float'),
       ((SELECT id FROM order_shipped), 'Order Total Without Discounts', 'order_total_without_discounts', 'float'),
       ((SELECT id FROM order_shipped), 'Postage Total', 'postage_total', 'float');
