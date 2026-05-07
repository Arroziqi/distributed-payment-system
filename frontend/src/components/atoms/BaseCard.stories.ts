import type { Meta, StoryObj } from '@storybook/vue3';
import BaseCard from './BaseCard.vue';
import BaseButton from './BaseButton.vue';

const meta = {
  title: 'Atoms/BaseCard',
  component: BaseCard,
  tags: ['autodocs'],
  argTypes: {
    title: { control: 'text' },
  },
} satisfies Meta<typeof BaseCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: (args) => ({
    components: { BaseCard },
    setup() {
      return { args };
    },
    template: `
      <BaseCard v-bind="args">
        <div>This is the card content.</div>
      </BaseCard>
    `,
  }),
  args: {
    title: 'Card Title',
  },
};

export const WithActions: Story = {
  render: (args) => ({
    components: { BaseCard, BaseButton },
    setup() {
      return { args };
    },
    template: `
      <BaseCard v-bind="args">
        <div>This is the card content.</div>
        <template #actions>
          <BaseButton label="Cancel" flat color="primary" />
          <BaseButton label="Save" color="primary" />
        </template>
      </BaseCard>
    `,
  }),
  args: {
    title: 'Card With Actions',
  },
};
