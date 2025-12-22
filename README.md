# Job Board API - Backend Project

**Duration**: 1 Week  
**Stack**: Go, Echo, PostgreSQL, JWT

---

## Overview

This is a RESTful API for a Job Board platform. Employers can post jobs, and job seekers can apply for them. The API provides authentication, CRUD operations for jobs and applications, and proper role-based access control.

---

## Technical Stack

- Language: Go 1.21+
- Framework: Echo v4
- Database: PostgreSQL
- Authentication: JWT
- Testing: Unit tests

---

## Features

### Authentication
- User registration (job_seeker or employer role)
- JWT-based login
- Protected routes based on role

### Employer Features
- Company profile management
- Job posting CRUD
- View applications for their jobs
- Update application status

### Job Seeker Features
- Profile management with skills
- Browse/search jobs
- Apply to jobs
- View own applications

### Public Features
- List published jobs with pagination
- Search and filter jobs
- View company profiles

### Application Workflow
- Job seekers apply with a cover letter
- Employers update status: `pending` → `reviewed` → `shortlisted` → `rejected` / `hired`
- Prevent duplicate applications

---

## API Standards

- RESTful conventions
- JSON responses with consistent structure
- Proper HTTP status codes
- Pagination support (page, per_page, total)

---

## Project Structure

