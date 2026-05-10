# Design System

This project utilizes a modern design system built on **Tailwind CSS** and **Shadcn-vue**, ensuring a highly customizable, accessible, and performant user interface.

## Core Principles
- **Utility-First**: Leverage Tailwind CSS for rapid styling without leaving the HTML.
- **Accessibility**: Use Shadcn-vue (built on Radix Vue) to ensure all components meet ARIA standards.
- **Consistency**: Centralized configuration in `tailwind.config.js`.

## Color Palette (Tailwind Tokens)
We use a curated HSL-based color system that supports light and dark modes seamlessly:

- **Primary**: Brand color for main actions and highlights.
- **Secondary**: For less prominent UI elements.
- **Destructive**: Used for dangerous actions (e.g., withdrawals, account deletion).
- **Muted**: For background elements and secondary text.
- **Accent**: Used for hover states and subtle highlights.

## Typography
- **Font Family**: Inter (sans-serif) - chosen for its clarity and modern feel.
- **Scale**: Defined using standard Tailwind text sizes (`text-sm`, `text-lg`, `text-2xl`, etc.).

## Components (Shadcn-vue)
Our UI is composed of highly reusable primitives located in `src/components/ui/`:
- **Buttons**: Variant-based (`default`, `outline`, `ghost`, `destructive`).
- **Cards**: For grouping dashboard information and transaction details.
- **Dialogs/Modals**: Accessible overlays for forms (e.g., Top-up, Transfer).
- **Toasts**: Auto-dismissing notifications for user feedback (powered by `vue-sonner`).
- **Progress Bars**: Used in modals and profile update feedback.

## Layout & Spacing
- **Grid**: 4px base unit, standardizing on Tailwind's spacing scale (`p-4`, `m-2`, `gap-6`).
- **Responsive**: Mobile-first design using Tailwind's breakpoints (`sm`, `md`, `lg`, `xl`).

## Design Tokens
Tokens are managed in `frontend/tailwind.config.js`, including custom animations and color extensions.
