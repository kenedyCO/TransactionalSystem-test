CREATE TABLE wallet (
    id uuid DEFAULT gen_random_uuid () primary key,
    actual int,
    frozen int
);

CREATE TABLE transactions (
    id uuid DEFAULT gen_random_uuid () primary key,
    amount int,
    status varchar(10),
    created_at timestamp DEFAULT now()
);



