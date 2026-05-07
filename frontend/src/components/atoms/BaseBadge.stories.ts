import type { Meta, StoryObj } from '@storybook/vue3';
import BaseBadge from './BaseBadge.vue';

const meta: Meta<typeof BaseBadge> = {
  title: 'Atoms/BaseBadge',
  component: BaseBadge,
};

export default meta;
type Story = StoryObj<typeof BaseBadge>;

export const Default: Story = {
  render: (args) => ({
    components: { BaseBadge },
    setup() {
      return { args };
    },
    template: '<BaseBadge v-bind="args">Badge</BaseBadge>',
  }),
  args: {
    variant: 'default',
  },
};
