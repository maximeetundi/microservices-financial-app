"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  ArrowLeftIcon,
  EyeIcon,
  EyeSlashIcon,
  KeyIcon,
} from "@heroicons/react/24/outline";

declare const process: { env: { NEXT_PUBLIC_API_URL?: string } };

const getApiUrl = () =>
  process.env.NEXT_PUBLIC_API_URL ||
  (typeof window !== "undefined" ? window.location.origin : "");

interface ProviderCredentials {
  [key: string]: string;
}

interface PlatformAccount {
  id: string;
  currency: string;
  balance: number;
  alias: string;
  name?: string;
  account_type?: string;
}

// Credentials field configuration per provider
const PROVIDER_CREDENTIALS: Record<
  string,
  { key: string; label: string; required: boolean; secret: boolean }[]
> = {
  cinetpay: [
    { key: "api_key", label: "API Key", required: true, secret: true },
    { key: "site_id", label: "Site ID", required: true, secret: false },
    { key: "secret_key", label: "Secret Key", required: false, secret: true },
  ],
  wave: [
    { key: "api_key", label: "API Key", required: true, secret: true },
    { key: "secret_key", label: "Secret Key", required: false, secret: true },
  ],
  mtn_money: [
    { key: "api_user", label: "API User", required: true, secret: false },
    { key: "api_key", label: "API Key", required: true, secret: true },
    {
      key: "subscription_key",
      label: "Subscription Key",
      required: true,
      secret: true,
    },
  ],
  orange_money: [
    { key: "client_id", label: "Client ID", required: true, secret: false },
    {
      key: "client_secret",
      label: "Client Secret",
      required: true,
      secret: true,
    },
    {
      key: "merchant_key",
      label: "Merchant Key",
      required: true,
      secret: true,
    },
  ],
  paypal: [
    { key: "client_id", label: "Client ID", required: true, secret: false },
    {
      key: "client_secret",
      label: "Client Secret",
      required: true,
      secret: true,
    },
    {
      key: "mode",
      label: "Mode (sandbox/live)",
      required: false,
      secret: false,
    },
  ],
  stripe: [
    { key: "api_key", label: "API Key", required: true, secret: true },
    {
      key: "public_key",
      label: "Public Key",
      required: false,
      secret: false,
    },
  ],
  flutterwave: [
    { key: "secret_key", label: "Secret Key", required: true, secret: true },
    { key: "public_key", label: "Public Key", required: true, secret: false },
    {
      key: "encryption_key",
      label: "Encryption Key",
      required: false,
      secret: true,
    },
  ],
  default: [
    { key: "api_key", label: "API Key", required: false, secret: true },
    { key: "secret_key", label: "Secret Key", required: false, secret: true },
    { key: "client_id", label: "Client ID", required: false, secret: false },
    {
      key: "client_secret",
      label: "Client Secret",
      required: false,
      secret: true,
    },
  ],
};

export default function NewAggregatorInstancePage() {
  const params = useParams();
  const router = useRouter();

  const providerId = params.id as string;

  const [providerName, setProviderName] = useState("");
  const [providerCode, setProviderCode] = useState("");
  const [hotWallets, setHotWallets] = useState<PlatformAccount[]>([]);
  const [selectedWalletIds, setSelectedWalletIds] = useState<string[]>([]);

  const [newInstance, setNewInstance] = useState({
    name: "",
    vault_secret_path: "",
    is_active: true,
    is_primary: false,
    is_global: false,
    is_test_mode: true,
    deposit_enabled: true,
    withdraw_enabled: true,
    priority: 50,
  });

  const [createCredentials, setCreateCredentials] = useState<ProviderCredentials>(
    {},
  );
  const [showCreateSecrets, setShowCreateSecrets] = useState<
    Record<string, boolean>
  >({});

  const getCredentialFields = () => {
    const code = providerCode.toLowerCase();
    return PROVIDER_CREDENTIALS[code] || PROVIDER_CREDENTIALS.default;
  };

  const vaultPathPreview = useMemo(() => {
    if (newInstance.vault_secret_path) return newInstance.vault_secret_path;
    const providerSlug = providerName.toLowerCase().replace(/\s+/g, "_");
    const instanceSlug =
      newInstance.name.toLowerCase().replace(/\s+/g, "_") || "instance_name";
    return `secret/aggregators/${providerSlug}/${instanceSlug}`;
  }, [newInstance.vault_secret_path, newInstance.name, providerName]);

  const toggleWalletSelection = (walletId: string) => {
    setSelectedWalletIds((prev: string[]) =>
      prev.includes(walletId)
        ? prev.filter((id: string) => id !== walletId)
        : [...prev, walletId],
    );
  };

  const toggleAllWallets = () => {
    if (selectedWalletIds.length === hotWallets.length) {
      setSelectedWalletIds([]);
    } else {
      setSelectedWalletIds(hotWallets.map((w: PlatformAccount) => w.id));
    }
  };

  const loadProviderDetails = async () => {
    try {
      const API_URL = getApiUrl();
      const token = localStorage.getItem("admin_token");
      const response = await fetch(
        `${API_URL}/api/v1/admin/payment-providers/${providerId}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (!response.ok) {
        throw new Error("Failed to load provider details");
      }

      const data = await response.json();
      const provider = data.provider || data;

      const name = String(provider?.name || provider?.display_name || "");
      setProviderName(name);

      const code = String(
        provider?.code || provider?.name || provider?.provider_code || "default",
      ).toLowerCase();
      setProviderCode(code);
    } catch (e) {
      console.error("Failed to load provider details", e);
    }
  };

  const loadHotWallets = async () => {
    try {
      const API_URL = getApiUrl();
      const token = localStorage.getItem("admin_token");
      const response = await fetch(
        `${API_URL}/api/v1/admin/platform/accounts?type=operations`,
        { headers: { Authorization: `Bearer ${token}` } },
      );
      if (response.ok) {
        const data = await response.json();
        setHotWallets(data.accounts || []);
      }
    } catch (e) {
      console.error("Failed to load hot wallets", e);
    }
  };

  useEffect(() => {
    if (providerId) {
      loadProviderDetails();
      loadHotWallets();
    }
  }, [providerId]);

  const handleCreateInstance = async (e: FormEvent) => {
    e.preventDefault();
    try {
      const API_URL = getApiUrl();
      const token = localStorage.getItem("admin_token");

      const path =
        newInstance.vault_secret_path ||
        `secret/aggregators/${providerName.toLowerCase().replace(/\s+/g, "_")}/${newInstance.name.toLowerCase().replace(/\s+/g, "_")}`;

      const walletsPayload = selectedWalletIds.map((id: string) => {
        const w = hotWallets.find((hw: PlatformAccount) => hw.id === id);
        return {
          hot_wallet_id: id,
          currency: w?.currency || "XOF",
        };
      });

      const credentialValues = Object.values(createCredentials) as string[];
      const hasRealCreateCredentials = credentialValues.some((value: string) => {
        const v = String(value || "");
        return v.trim() !== "" && !v.startsWith("****");
      });

      const requestBody: any = {
        ...newInstance,
        vault_secret_path: path,
        wallets: walletsPayload,
      };

      if (hasRealCreateCredentials) {
        const filteredCreds: ProviderCredentials = {};
        Object.entries(createCredentials).forEach(([key, value]) => {
          const v = String(value || "");
          if (v && v.trim() !== "" && !v.startsWith("****")) {
            filteredCreds[key] = v;
          }
        });
        requestBody.credentials = filteredCreds;
      }

      const response = await fetch(
        `${API_URL}/api/v1/admin/payment-providers/${providerId}/instances`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody),
        },
      );

      if (response.ok) {
        router.push(`/dashboard/aggregators/${providerId}`);
      } else {
        const errData = await response.json();
        alert(`Failed to create instance: ${errData.error || "Unknown error"}`);
      }
    } catch (e) {
      console.error(e);
      alert("Failed to create instance: Network error");
    }
  };

  return (
    <div className="space-y-6 animate-fadeIn">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.back()}
            className="p-2 rounded-lg hover:bg-gray-100 text-gray-500"
          >
            <ArrowLeftIcon className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Nouvelle Instance</h1>
            <p className="text-gray-500 text-sm">
              {providerName ? `Provider: ${providerName}` : "Chargement..."}
            </p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
        <form onSubmit={handleCreateInstance} className="space-y-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Nom (ex: Instance 1)
            </label>
            <input
              type="text"
              required
              className="input w-full"
              value={newInstance.name}
              onChange={(e) => setNewInstance({ ...newInstance, name: e.target.value })}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Chemin Vault
            </label>
            <input
              type="text"
              className="input w-full font-mono text-sm"
              placeholder="Généré automatiquement si vide"
              value={newInstance.vault_secret_path}
              onChange={(e) =>
                setNewInstance({ ...newInstance, vault_secret_path: e.target.value })
              }
            />
            <div className="mt-2 p-2 bg-gray-100 rounded-lg border border-gray-200">
              <p className="text-xs text-gray-500 mb-1">Chemin Vault utilisé:</p>
              <code className="text-xs font-mono text-indigo-600 break-all">
                {vaultPathPreview}
              </code>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="col-span-2">
              <div className="text-sm font-medium text-gray-700 mb-2">Mode</div>
              <div className="flex gap-3">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    name="instance_mode"
                    checked={newInstance.is_test_mode === true}
                    onChange={() => setNewInstance({ ...newInstance, is_test_mode: true })}
                  />
                  <span className="text-sm">Sandbox</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    name="instance_mode"
                    checked={newInstance.is_test_mode === false}
                    onChange={() => setNewInstance({ ...newInstance, is_test_mode: false })}
                  />
                  <span className="text-sm">Live</span>
                </label>
              </div>
            </div>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox"
                checked={newInstance.is_active}
                onChange={(e) => setNewInstance({ ...newInstance, is_active: e.target.checked })}
              />
              <span className="text-sm">Active</span>
            </label>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox"
                checked={newInstance.deposit_enabled}
                onChange={(e) =>
                  setNewInstance({ ...newInstance, deposit_enabled: e.target.checked })
                }
              />
              <span className="text-sm">Dépôts</span>
            </label>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox"
                checked={newInstance.withdraw_enabled}
                onChange={(e) =>
                  setNewInstance({ ...newInstance, withdraw_enabled: e.target.checked })
                }
              />
              <span className="text-sm">Retraits</span>
            </label>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox"
                checked={newInstance.is_primary}
                onChange={(e) => setNewInstance({ ...newInstance, is_primary: e.target.checked })}
              />
              <span className="text-sm">Principal</span>
            </label>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox"
                checked={newInstance.is_global}
                onChange={(e) => setNewInstance({ ...newInstance, is_global: e.target.checked })}
              />
              <span className="text-sm">🌍 Global (tous pays)</span>
            </label>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Priorité</label>
            <input
              type="number"
              className="input w-full"
              value={newInstance.priority}
              onChange={(e) =>
                setNewInstance({ ...newInstance, priority: parseInt(e.target.value) })
              }
            />
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <label className="block text-sm font-medium text-gray-700">
                Hot Wallets Associés
              </label>
              <button
                type="button"
                onClick={toggleAllWallets}
                className="text-xs text-blue-600 hover:text-blue-800"
              >
                {selectedWalletIds.length === hotWallets.length
                  ? "Tout désélectionner"
                  : "Tout sélectionner"}
              </button>
            </div>
            <div className="border rounded-md p-3 max-h-48 overflow-y-auto space-y-2 bg-gray-50">
              {hotWallets.length === 0 ? (
                <p className="text-gray-500 text-sm">Aucun hot wallet disponible.</p>
              ) : (
                hotWallets.map((wallet) => (
                  <label
                    key={wallet.id}
                    className="flex items-center space-x-2 cursor-pointer p-1 hover:bg-gray-100 rounded"
                  >
                    <input
                      type="checkbox"
                      className="checkbox checkbox-xs"
                      checked={selectedWalletIds.includes(wallet.id)}
                      onChange={() => toggleWalletSelection(wallet.id)}
                    />
                    <span className="text-sm">
                      <span className="font-semibold">{wallet.currency}</span> - {wallet.name}
                      <span className="text-gray-500 text-xs ml-1">
                        (Solde: {wallet.balance})
                      </span>
                    </span>
                  </label>
                ))
              )}
            </div>
            <p className="text-xs text-gray-500 mt-1">
              Associer des wallets pour permettre les dépôts/retraits immédiats.
            </p>
          </div>

          <div className="border-t border-gray-100 pt-4">
            <div className="flex items-center gap-2 mb-2">
              <KeyIcon className="w-4 h-4 text-gray-400" />
              <div className="text-sm font-medium text-gray-700">Credentials (optionnel)</div>
            </div>
            <div className="space-y-3">
              {getCredentialFields().map((field) => (
                <div key={field.key}>
                  <label className="block text-xs font-medium text-gray-600 mb-1">
                    {field.label}
                    {field.required ? " *" : ""}
                  </label>
                  <div className="relative">
                    <input
                      type={field.secret && !showCreateSecrets[field.key] ? "password" : "text"}
                      className="input w-full pr-10"
                      value={createCredentials[field.key] || ""}
                      onChange={(e) =>
                        setCreateCredentials({
                          ...createCredentials,
                          [field.key]: e.target.value,
                        })
                      }
                    />
                    {field.secret && (
                      <button
                        type="button"
                        onClick={() =>
                          setShowCreateSecrets({
                            ...showCreateSecrets,
                            [field.key]: !showCreateSecrets[field.key],
                          })
                        }
                        className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                      >
                        {showCreateSecrets[field.key] ? (
                          <EyeSlashIcon className="w-4 h-4" />
                        ) : (
                          <EyeIcon className="w-4 h-4" />
                        )}
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={() => router.back()}
              className="btn-secondary"
            >
              Annuler
            </button>
            <button type="submit" className="btn-primary">
              Créer
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
