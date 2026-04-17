-- Mock Event Types
INSERT INTO prism.event_types (name, event_key, version, description) VALUES
('Page View', 'page_view', 1, 'Standard page view event.'),
('Button Click', 'button_click', 1, 'Standard button click event.'),
('Purchase', 'purchase', 1, 'Event triggered on a successful purchase.'),
('Add to Cart', 'add_to_cart', 1, 'Event triggered when an item is added to cart.')
ON CONFLICT (name) DO NOTHING;

-- Get IDs of event types and insert fields
DO $$
DECLARE
    page_view_id UUID;
    button_click_id UUID;
    purchase_id UUID;
    add_to_cart_id UUID;
BEGIN
    SELECT id INTO page_view_id FROM prism.event_types WHERE event_key = 'page_view';
    SELECT id INTO button_click_id FROM prism.event_types WHERE event_key = 'button_click';
    SELECT id INTO purchase_id FROM prism.event_types WHERE event_key = 'purchase';
    SELECT id INTO add_to_cart_id FROM prism.event_types WHERE event_key = 'add_to_cart';

    -- Event Fields for Page View
    INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type) VALUES
    (page_view_id, 'Page Path', 'page_path', 'string'),
    (page_view_id, 'User Agent', 'user_agent', 'string')
    ON CONFLICT (event_type_id, field_key) DO NOTHING;

    -- Event Fields for Button Click
    INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type) VALUES
    (button_click_id, 'Button ID', 'button_id', 'string'),
    (button_click_id, 'Page Path', 'page_path', 'string')
    ON CONFLICT (event_type_id, field_key) DO NOTHING;

    -- Event Fields for Purchase
    INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type) VALUES
    (purchase_id, 'Amount', 'amount', 'float'),
    (purchase_id, 'Currency', 'currency', 'string')
    ON CONFLICT (event_type_id, field_key) DO NOTHING;

    -- Event Fields for Add to Cart
    INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type) VALUES
    (add_to_cart_id, 'Product ID', 'product_id', 'string')
    ON CONFLICT (event_type_id, field_key) DO NOTHING;
END $$;

-- Simple Metrics
INSERT INTO prism.metrics (name, metric_key, description, metric_type, analysis_unit) VALUES
('Purchase Count', 'purchase_count', 'Total number of purchases.', 'simple', 'user'),
('Average Purchase Amount', 'avg_purchase_amount', 'Average amount of a purchase.', 'simple', 'user'),
('Total Purchase Value', 'total_purchase_value', 'Total value of all purchases.', 'simple', 'user'),
('Page View Count', 'page_view_count', 'Total number of page views.', 'simple', 'user'),
('Button Click Count', 'button_click_count', 'Total number of button clicks.', 'simple', 'user'),
('Add to Cart Count', 'add_to_cart_count', 'Total number of items added to cart.', 'simple', 'user')
ON CONFLICT (name) DO NOTHING;

-- Metric Components
DO $$
DECLARE
    purchase_count_id UUID;
    avg_purchase_amount_id UUID;
    total_purchase_value_id UUID;
    page_view_count_id UUID;
    button_click_count_id UUID;
    add_to_cart_count_id UUID;
    
    purchase_type_id UUID;
    page_view_type_id UUID;
    button_click_type_id UUID;
    add_to_cart_type_id UUID;
    
    purchase_amount_field_id UUID;
BEGIN
    SELECT id INTO purchase_count_id FROM prism.metrics WHERE metric_key = 'purchase_count';
    SELECT id INTO avg_purchase_amount_id FROM prism.metrics WHERE metric_key = 'avg_purchase_amount';
    SELECT id INTO total_purchase_value_id FROM prism.metrics WHERE metric_key = 'total_purchase_value';
    SELECT id INTO page_view_count_id FROM prism.metrics WHERE metric_key = 'page_view_count';
    SELECT id INTO button_click_count_id FROM prism.metrics WHERE metric_key = 'button_click_count';
    SELECT id INTO add_to_cart_count_id FROM prism.metrics WHERE metric_key = 'add_to_cart_count';

    SELECT id INTO purchase_type_id FROM prism.event_types WHERE event_key = 'purchase';
    SELECT id INTO page_view_type_id FROM prism.event_types WHERE event_key = 'page_view';
    SELECT id INTO button_click_type_id FROM prism.event_types WHERE event_key = 'button_click';
    SELECT id INTO add_to_cart_type_id FROM prism.event_types WHERE event_key = 'add_to_cart';

    SELECT id INTO purchase_amount_field_id FROM prism.event_fields WHERE event_type_id = purchase_type_id AND field_key = 'amount';

    -- Components for Purchase Count
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation) VALUES
    (purchase_count_id, 'base_event', purchase_type_id, 'COUNT');

    -- Components for Average Purchase Amount
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation, agg_field_id) VALUES
    (avg_purchase_amount_id, 'base_event', purchase_type_id, 'AVG', purchase_amount_field_id);

    -- Components for Total Purchase Value
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation, agg_field_id) VALUES
    (total_purchase_value_id, 'base_event', purchase_type_id, 'SUM', purchase_amount_field_id);

    -- Components for Page View Count
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation) VALUES
    (page_view_count_id, 'base_event', page_view_type_id, 'COUNT');

    -- Components for Button Click Count
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation) VALUES
    (button_click_count_id, 'base_event', button_click_type_id, 'COUNT');

    -- Components for Add to Cart Count
    INSERT INTO prism.metric_components (metric_id, role, event_type_id, agg_operation) VALUES
    (add_to_cart_count_id, 'base_event', add_to_cart_type_id, 'COUNT');
END $$;
