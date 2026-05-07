import { computed } from 'vue';
import { useGetTransactions, usePostTransactionsPayments } from '@/api/generated/transactions/transactions';
import { toast } from 'vue-sonner';

export const useTransactions = (params?: any) => {
  const transactionsQuery = useGetTransactions(params, {
    query: {
      refetchInterval: 30000,
    }
  });

  const paymentMutation = usePostTransactionsPayments();

  const transactions = computed(() => {
    const data = transactionsQuery.data.value as any;
    return data?.transactions || [];
  });

  const processPayment = async (data: any, idempotencyKey: string) => {
    try {
      await paymentMutation.mutateAsync({
        data,
        // @ts-ignore
        headers: {
          'Idempotency-Key': idempotencyKey
        }
      });
      toast.success('Payment processed successfully');
      transactionsQuery.refetch();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Payment failed');
    }
  };

  return {
    transactions,
    loading: transactionsQuery.isLoading,
    processPayment,
    isProcessing: paymentMutation.isPending,
  };
};
