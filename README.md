# FiberNova
FiberNova is a full-stack, Laravel-inspired web framework built on Go Fiber, designed to combine blazing-fast performance with elegant, structured development for modern cloud-native applications.

## Getting Started
Make command for:



    Command to initial FiberNova app

    Create Model

    Create Middleware
Create controller
Create route

Here are the proposed CLI commands for initializing and developing a FiberNova application:

### Initialize FiberNova App
```bash
fibernova new myapp
```
Creates a new FiberNova project with full directory structure, configuration files, and dependencies installed.

### Create Model
```bash
fibernova make:model User
```
Generates a new model file `User.go` under the `app/Models` directory with basic struct and database mapping.

### Create Controller
```bash
fibernova make:controller UserController
```
Generates a new controller `UserController.go` in `app/Http/Controllers` with resource method stubs (index, show, store, update, destroy).

### Create Middleware
```bash
fibernova make:middleware AuthMiddleware
```
Creates a middleware file `AuthMiddleware.go` in `app/Http/Middleware` with `Handle` method ready for logic implementation.

### Create Route
Routes are defined in `routes/web.go` or `routes/api.go`. To generate a resource route scaffold:
```bash
fibernova make:route resource /users UserController
```
This appends RESTful route bindings for `/users` to the router, mapped to `UserController` methods.

The CLI follows Laravel-inspired conventions while adapting to Go’s package structure and Fiber’s routing system, ensuring a smooth developer experience.

## Technology Stack
[Fiber](https://gofiber.io/) + [GORM](https://gorm.io/index.html) + [toml](https://toml.io/en/) + [Nuxt.js](https://nuxt.com/)
