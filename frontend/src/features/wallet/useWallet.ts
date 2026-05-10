import { computed } from 'vue';
import { useGetWalletsUserIDBalance, usePostWalletTopups, usePostWalletTransfers, usePostWallets } from '@/api/generated/wallet/wallet';
import { toast } from 'vue-sonner';
import { useAuthStore } from '@/stores/auth.store';

export const useWallet = () => {
  const authStore = useAuthStore();
  const userId = authStore.currentUser?.id || '';

  const balanceQuery = useGetWalletsUserIDBalance(userId, {
    query: {
      enabled: !!userId,
      refetchInterval: 10000, // Poll every 10s
    }
  });

  const topupMutation = usePostWalletTopups();
  const createWalletMutation = usePostWallets();
  const transferMutation = usePostWalletTransfers();

  const wallet = computed(() => {
    const data = balanceQuery.data.value as any;
    return data?.balance;
  });

  const balance = computed(() => {
    return wallet.value?.available_balance || 0;
  });

  const lockedBalance = computed(() => {
    return wallet.value?.locked_balance || 0;
  });

  const currency = computed(() => {
    return wallet.value?.currency || 'USD';
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

  const createWallet = async (currency?: string) => {
    try {
      await createWalletMutation.mutateAsync({
        data: { user_id: userId, currency }
      });
      toast.success('Wallet created successfully');
      balanceQuery.refetch();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Create wallet failed');
      throw error;
    }
  }

  return {
    wallet,
    balance,
    lockedBalance,
    currency,
    loading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
    error: balanceQuery.error,
    topup,
    transfer,
    createWallet,
    isTopupLoading: topupMutation.isPending,
    isTransferLoading: transferMutation.isPending,
    isCreateWalletLoading: createWalletMutation.isPending,
    refetch: balanceQuery.refetch,
  };
};
