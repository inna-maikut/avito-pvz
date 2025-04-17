CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email text not null,
    password text not null,
    user_role SMALLINT not null,
    create_time timestamp with time zone default now()
);
CREATE unique INDEX users__email on users (email);

CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID REFERENCES pvz(id),
    status SMALLINT NOT NULL,
    recepted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX receptions__recepted_at ON receptions(recepted_at);
CREATE INDEX receptions__pvz_id_status ON receptions(pvz_id, status);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID REFERENCES receptions(id),
    category SMALLINT NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()   
);

CREATE INDEX products__reception_id_added_at 
    ON products(reception_id, added_at);


