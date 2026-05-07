import type { Meta, StoryObj } from '@storybook/vue3';
import BaseButton from './BaseButton.vue';

const meta: Meta<typeof BaseButton> = {
  title: 'Atoms/BaseButton',
  component: BaseButton,
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'destructive', 'outline', 'secondary', 'ghost', 'link'],
    },
    size: {
      control: 'select',
      options: ['default', 'sm', 'lg', 'icon'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof BaseButton>;

export const Default: Story = {
  render: (args) => ({
    components: { BaseButton },
    setup() {
      return { args };
    },
    template: '<BaseButton v-bind="args">Button</BaseButton>',
  }),
  args: {
    variant: 'default',
    size: 'default',
  },
};

export const Loading: Story = {
  ...Default,
  args: {
    ...Default.args,
    loading: true,
  },
};

export const Destructive: Story = {
  ...Default,
  args: {
    ...Default.args,
    variant: 'destructive',
  },
};
