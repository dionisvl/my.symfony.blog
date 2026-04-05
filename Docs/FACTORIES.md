# Demo Data & Factories

## Quick Setup

Run this command to populate the database with demo data:

```bash
php artisan migrate:fresh --seed
# Or
php artisan db:seed
```

## Available Factories

### Blog Models

- **User**: Creates users with admin/regular roles
- **Post**: Creates blog posts with content, categories, and tags
- **Category**: Creates post categories
- **Tag**: Creates post tags

### Manual Factory Usage

Enter tinker mode:

```bash
php artisan tinker
```

Then run individual factories:

```php
User::factory()->create();
User::factory()->create(['is_admin' => true]);

Category::factory()->count(5)->create();
Tag::factory()->count(10)->create();

Post::factory()->count(10)->create();

Post::factory()->count(3)->create(['is_featured' => 1]);
```

## Demo Data Structure

The `BlogSeeder` creates:

- **Admin user**: `admin@blog.local` / `admin123`
- **Regular user**: `user@blog.local` / `user123`
- **5 categories** with realistic names
- **10 tags** for content organization
- **20 regular posts** with random content
- **5 featured posts** for homepage highlights

## Custom Seeder Commands

### Full Reset

```bash
php artisan migrate:fresh --seed
```

### Specific Seeder

```bash
php artisan db:seed --class=BlogSeeder
```

## Factory Relationships

Posts are automatically linked to:

- Random categories
- Random tags (1-3 per post)
- Admin user as author
