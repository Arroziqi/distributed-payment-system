# Storybook Guide

Storybook is an open source tool for UI component development and testing. It allows developers to create components independently and showcase components interactively in an isolated development environment.

## Accessing Storybook
- Locally via Docker: Navigate to `http://localhost:6006`
- The `storybook` container automatically spins up alongside the other services when running `docker compose up`.

## Writing Stories
Stories are defined in the same directory as the component they document, following the `.stories.ts` naming convention.

### Requirements for each Story:
1. **Default State**: How the component looks normally.
2. **Loading State**: How the component indicates background activity.
3. **Disabled State**: How the component looks when interactions are prevented.
4. **Error State**: (If applicable) How the component displays validation errors.

### Example Structure:
```typescript
import type { Meta, StoryObj } from '@storybook/vue3';
import MyComponent from './MyComponent.vue';

const meta = {
  title: 'Atoms/MyComponent',
  component: MyComponent,
  tags: ['autodocs'],
} satisfies Meta<typeof MyComponent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = { args: { ... } };
```
