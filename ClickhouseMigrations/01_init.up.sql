CREATE TABLE IF NOT EXISTS redirects
(
    id UInt32, 
    original String, 
    short String, 
    user_ip String, 
    os String, 
    created_at DateTime DEFAULT now()
) 
ENGINE = MergeTree()
ORDER BY id;
