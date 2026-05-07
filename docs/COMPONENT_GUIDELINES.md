# Component Guidelines

## Rules
1. **Clear Naming:** Component names should be multi-word (e.g., `BaseButton`, `LoginForm`) to prevent conflicts with native HTML elements.
2. **Prop Validation:** Use TypeScript interfaces with `defineProps` to enforce strict typing on all component inputs.
3. **Event Emitting:** Use `defineEmits` to declare custom events. Do not mutate props directly.
4. **Slots:** Use slots for content injection to maximize reusability, especially for structural components like Cards and Modals.
5. **No Direct API Calls in Atoms/Molecules:** Data fetching should happen at the Page level, and data should be passed down via props.

## Props vs State
- **Props:** Data passed from a parent component. Immutable inside the child component.
- **State:** Local data managed within the component (using `ref` or `reactive`).

## Styling
- Component-specific styles should be scoped `<style scoped lang="sass">`.
- Use Quasar utility classes (e.g., `q-mt-md`, `text-primary`) whenever possible instead of writing custom CSS.
