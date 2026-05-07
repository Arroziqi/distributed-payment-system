import type { Meta, StoryObj } from '@storybook/vue3';
import LoginForm from './LoginForm.vue';

const meta = {
  title: 'Molecules/LoginForm',
  component: LoginForm,
  tags: ['autodocs'],
  argTypes: {
    loading: { control: 'boolean' },
  },
} satisfies Meta<typeof LoginForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    loading: false,
  },
};

export const Loading: Story = {
  args: {
    loading: true,
  },
};
