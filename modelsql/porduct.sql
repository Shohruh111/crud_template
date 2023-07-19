
CREATE TABLE "product" (
    "id" UUID PRIMARY KEY,
    "title" VARCHAR NOT NULL,
    "price" SERIAL NOT NULL,
    "category" UUID REFERENCES "category"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);
