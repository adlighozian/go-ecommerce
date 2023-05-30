CREATE TABLE "api_managements" (
  "id" int PRIMARY KEY,
  "api_name" text UNIQUE NOT NULL,
  "service_name" text NOT NULL,
  "endpoint_url" text NOT NULL,
  "hashed_endpoint_url" text UNIQUE NOT NULL,
  "is_available" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "splash_images" (
  "id" int PRIMARY KEY,
  "image_url" varchar(255) UNIQUE NOT NULL,
  "is_active" boolean NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "languages" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "username" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "role" varchar(255) NOT NULL,
  "full_name" varchar(255),
  "age" int,
  "image_url" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "user_settings" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "notification" boolean,
  "dark_mode" boolean,
  "language_id" int,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "addresses" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "street" varchar(255),
  "city" varchar(255),
  "state" varchar(255),
  "country_id" int,
  "zipcode" varchar(255),
  "phone_number" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "countries" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "calling_code" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "stores" (
  "id" int PRIMARY KEY,
  "address_id" int,
  "name" varchar(255),
  "description" varchar(255),
  "image_url" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "store_id" int,
  "category_id" int,
  "size_id" int,
  "color_id" int,
  "name" varchar(255),
  "subtitle" varchar(255),
  "description" text,
  "unit_price" float,
  "status" boolean,
  "stock" int,
  "SKU" varchar(255),
  "weight" float,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "product_sizes" (
  "id" int PRIMARY KEY,
  "size" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "product_colors" (
  "id" int PRIMARY KEY,
  "color" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "product_images" (
  "id" int PRIMARY KEY,
  "product_id" int,
  "image_url" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "wishlists" (
  "id" serial not null PRIMARY KEY,
  "user_id" int,
  "product_id" int,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "carts" (
  "id" serial not null PRIMARY KEY,
  "user_id" int,
  "product_id" int,
  "quantity" int,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "payment_method_id" int,
  "shipping_id" int,
  "total_price" float,
  "status" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" int PRIMARY KEY,
  "order_id" int,
  "product_id" int,
  "quantity" int,
  "unit_price" float,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "payment_methods" (
  "id" int PRIMARY KEY,
  "payment_gateway_id" text NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "shippings" (
  "id" int PRIMARY KEY,
  "receipt_number" text
);

CREATE TABLE "payment_logs" (
  "id" int PRIMARY KEY,
  "order_id" int,
  "payment_method_id" int,
  "payment_status" varchar(255),
  "transaction_id" varchar(255),
  "transaction_time" timestamp,
  "transaction_status" varchar(255),
  "gross_amount" float,
  "debug_info" json,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "reviews" (
  "id" serial not null PRIMARY KEY,
  "user_id" int,
  "product_id" int,
  "rating" int,
  "review_text" text,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "review_images" (
  "id" int PRIMARY KEY,
  "review_id" int,
  "image_url" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "promotions" (
  "id" int PRIMARY KEY,
  "store_id" int,
  "product_id" int,
  "category_id" int,
  "discount_type" varchar(255),
  "discount_value" float,
  "start_date" timestamp,
  "end_date" timestamp,
  "is_active" boolean NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "promotion_codes" (
  "id" int PRIMARY KEY,
  "promotion_id" int,
  "code" varchar(255),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "promotions"."start_date" IS 'keterangan tanggal voucher abis';

ALTER TABLE "user_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_settings" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "addresses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "addresses" ADD FOREIGN KEY ("country_id") REFERENCES "countries" ("id");

ALTER TABLE "stores" ADD FOREIGN KEY ("address_id") REFERENCES "addresses" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("size_id") REFERENCES "product_sizes" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("color_id") REFERENCES "product_colors" ("id");

ALTER TABLE "product_images" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("payment_method_id") REFERENCES "payment_methods" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("shipping_id") REFERENCES "shippings" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "payment_logs" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "payment_logs" ADD FOREIGN KEY ("payment_method_id") REFERENCES "payment_methods" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "review_images" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id");

ALTER TABLE "promotions" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "promotions" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "promotions" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "promotion_codes" ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id");
