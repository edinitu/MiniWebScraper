set search_path to 'go_prj';
set schema 'public';

drop table if exists shops;
create table shops(
	shop_id serial not null,
	shop_name text,
	constraint shops_pk primary key (shop_id),
    constraint shops_name_uk unique (shop_name)
);

drop table if exists products;
create table products(
	product_id serial not null,
	shop_id bigint,
	product_name text,
	price text,
	IsPromotion bool,
	OriginalPrice text,
	Quantity text,
	constraint product_pk primary key (product_id),
    constraint product_name_uk unique (product_name),
	constraint shop_id_fk foreign key (shop_id) 
		references shops(shop_id)
);

drop table if exists patterns;
create table patterns(
	pattern_id serial not null,
	pattern_value text,
	shop_id bigint,
	constraint pattern_pk primary key (pattern_id),
	constraint pattern_value_uk unique (pattern_value),
	constraint shop_id_pattern_fk foreign key (shop_id) 
		references shops(shop_id)
);

drop table if exists categories;
create table categories(
	category_id serial not null,
	category_name text,
	shop_id bigint,
	constraint category_pk primary key (category_id),
	constraint category_value_uk unique (category_name),
	constraint shop_id_category_fk foreign key (shop_id) 
		references shops(shop_id)
);







