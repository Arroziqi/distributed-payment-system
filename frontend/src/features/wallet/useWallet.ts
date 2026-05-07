import { computed } from 'vue';
import { useGetWalletsUserIDBalance, usePostWalletTopups, usePostWalletTransfers } from '@/api/generated/wallet/wallet';
import { toast } from 'vue-sonner';

export const useWallet = () => {
  const userId = localStorage.getItem('user_id') || '';

  const balanceQuery = useGetWalletsUserIDBalance(userId, {
    query: {
      enabled: !!userId,
      refetchInterval: 10000, // Poll every 10s
    }
  });

  const topupMutation = usePostWalletTopups();
  const transferMutation = usePostWalletTransfers();

  const balance = computed(() => {
    const data = balanceQuery.data.value as any;
    return data?.balance || 0;
  });

  const topup = async (amount: number) => {
    try {
      await topupMutation.mutateAsync({
        data: {
          amount,
          user_id: userId
        }
      });
      toast.success('Top up successful');
      balanceQuery.refetch();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Top up failed');
    }
  };

  const transfer = async (toUserId: string, amount: number) => {
    try {
      await transferMutation.mutateAsync({
        data: {
          from_user_id: userId,
          to_user_id: toUserId,
          amount
        }
      });
      toast.success('Transfer successful');
      balanceQuery.refetch();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Transfer failed');
    }
  };

  return {
    balance,
    loading: balanceQuery.isLoading,
    topup,
    transfer,
    isTopupLoading: topupMutation.isPending,
    isTransferLoading: transferMutation.isPending,
  };
};
