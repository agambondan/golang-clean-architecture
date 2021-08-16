package main

import "fmt"

func main() {
	fmt.Println("CREATE TABLE IF NOT EXISTS users (uuid uuid PRIMARY KEY, first_name VARCHAR(55) not null, last_name VARCHAR(55) not null, email VARCHAR(55) unique not null, phone_number VARCHAR(15) not null, username VARCHAR(55) unique not null, password VARCHAR(255) not null, image VARCHAR(255), image_url VARCHAR(255), thumbnail_url VARCHAR(255), role_id int not null, instagram VARCHAR(255), facebook VARCHAR(255), twitter VARCHAR(255), linkedin VARCHAR(255), created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	fmt.Println("CREATE TABLE IF NOT EXISTS roles (id serial PRIMARY KEY, name VARCHAR(15) not null unique, created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	fmt.Println("CREATE TABLE IF NOT EXISTS posts (id serial PRIMARY KEY, title VARCHAR(100) not null unique, description text not null, image VARCHAR(255), image_url VARCHAR(255), thumbnail_url VARCHAR(255), user_uuid uuid not null, created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	fmt.Println("CREATE TABLE IF NOT EXISTS categories (id serial PRIMARY KEY, name VARCHAR(35) not null unique, image VARCHAR(255), image_url VARCHAR(255), thumbnail_url VARCHAR(255), created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	fmt.Println("CREATE TABLE IF NOT EXISTS post_categories (post_id serial not null, category_id int not null)")
}
