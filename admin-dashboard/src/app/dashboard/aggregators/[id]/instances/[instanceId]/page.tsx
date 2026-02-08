"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  ArrowLeftIcon,
  CheckCircleIcon,
  EyeIcon,
  EyeSlashIcon,
  KeyIcon,
  XCircleIcon,
} from "@heroicons/react/24/outline";

type CredentialsSource = "db" | "vault";

type ProviderCredentials = Record<string, string>;

interface ProviderDetails {
  id: string;
  code?: string;
  name?: string;
}

interface ProviderInstance {
  id: string;
  provider_id: string;
  name: string;
  vault_secret_path?: string;
  hot_wallet_id?: string;
  is_active: boolean;
  is_primary: boolean;
  is_global: boolean;
  is_paused: boolean;
  priority: number;
  request_count: number;
  health_status: string;
}

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
      key: "business_currencies",
      label: "Devises PayPal Business (ex: USD,EUR,GBP)",
      required: false,
      secret: false,
    },
    {
      key: "mode",
      label: "Mode (sandbox/live)",
      required: false,
      secret: false,
    },
    {
      key: "base_url",
      label: "Base URL (optional)",
      required: false,
      secret: false,
    },
  ],
  stripe: [
    {
      key: "api_key",
      label: "Secret Key (sk_...)",
      required: true,
      secret: true,
    },
    {
      key: "public_key",
      label: "Publishable Key (pk_...)",
      required: true,
      secret: false,
    },
    {
      key: "webhook_secret",
      label: "Webhook Secret",
      required: false,
      secret: true,
    },
  ],
  flutterwave: [
    { key: "public_key", label: "Public Key", required: true, secret: false },
    { key: "secret_key", label: "Secret Key", required: true, secret: true },
    {
      key: "encryption_key",
      label: "Encryption Key",
      required: true,
      secret: true,
    },
  ],
  paystack: [
    { key: "public_key", label: "Public Key", required: true, secret: false },
    { key: "secret_key", label: "Secret Key", required: true, secret: true },
  ],
  default: [
    { key: "api_key", label: "API Key", required: true, secret: true },
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

export default function AggregatorInstanceEditPage() {
  const router = useRouter();
  const params = useParams();

  const providerId = params?.id as string;
  const instanceId = params?.instanceId as string;

  const [providerDetails, setProviderDetails] = useState<ProviderDetails | null>(
    null,
  );
  const [providerCode, setProviderCode] = useState<string>("default");

  const [instance, setInstance] = useState<ProviderInstance | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string>("");

  const [source, setSource] = useState<CredentialsSource>("db");
  const [loadingCredentials, setLoadingCredentials] = useState(false);
  const [savingCredentials, setSavingCredentials] = useState(false);

  const [credentials, setCredentials] = useState<ProviderCredentials>({});
  const [showSecret, setShowSecret] = useState<Record<string, boolean>>({});

  const fields = useMemo(() => {
    return PROVIDER_CREDENTIALS[providerCode] || PROVIDER_CREDENTIALS.default;
  }, [providerCode]);

  const API_URL =
    process.env.NEXT_PUBLIC_API_URL || "http://localhost:8088";

  const getToken = () => localStorage.getItem("admin_token");

  const loadProvider = async () => {
    const resp = await fetch(
      `${API_URL}/api/v1/admin/payment-providers/${providerId}`,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
          "Content-Type": "application/json",
        },
      },
    );

    if (!resp.ok) {
      throw new Error("Failed to load provider details");
    }

    const data = await resp.json();
    const p = (data.provider || data) as ProviderDetails;
    setProviderDetails(p);
    const code = String((p as any)?.code || (p as any)?.name || "default").toLowerCase();
    setProviderCode(code || "default");
  };

  const loadInstance = async () => {
    const resp = await fetch(
      `${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}`,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
          "Content-Type": "application/json",
        },
      },
    );

    if (!resp.ok) {
      throw new Error("Failed to load instance");
    }

    const data = await resp.json();
    const inst = (data.instance || data) as ProviderInstance;
    setInstance(inst);
  };

  const loadCredentials = async (nextSource: CredentialsSource) => {
    setLoadingCredentials(true);
    try {
      const resp = await fetch(
        `${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}/credentials?source=${nextSource}`,
        {
          headers: {
            Authorization: `Bearer ${getToken()}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (!resp.ok) {
        setCredentials({});
        return;
      }

      const data = await resp.json();
      setCredentials((data.credentials || {}) as ProviderCredentials);
    } finally {
      setLoadingCredentials(false);
    }
  };

  const saveCredentials = async () => {
    setSavingCredentials(true);
    try {
      const filteredCreds: ProviderCredentials = {};
      Object.entries(credentials).forEach(([k, v]) => {
        const vv = String(v ?? "");
        if (vv && vv.trim() !== "" && !vv.startsWith("****")) {
          filteredCreds[k] = vv;
        }
      });

      const resp = await fetch(
        `${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}/credentials?source=${source}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${getToken()}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(filteredCreds),
        },
      );

      if (!resp.ok) {
        let errMsg = "Impossible de sauvegarder";
        try {
          const errData = await resp.json();
          errMsg = errData.error || errMsg;
        } catch {
          // ignore
        }
        alert(`Erreur: ${errMsg}`);
        return;
      }

      alert("✅ Credentials sauvegardés avec succès!");
      await loadCredentials(source);
    } finally {
      setSavingCredentials(false);
    }
  };

  useEffect(() => {
    let cancelled = false;
    const run = async () => {
      setLoading(true);
      setLoadError("");
      try {
        await Promise.all([loadProvider(), loadInstance()]);
        if (!cancelled) {
          await loadCredentials(source);
        }
      } catch (e: any) {
        if (!cancelled) {
          setLoadError(e?.message || "Erreur de chargement");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    if (providerId && instanceId) {
      run();
    }

    return () => {
      cancelled = true;
    };
  }, [providerId, instanceId]);

  const onChangeField = (key: string, value: string) => {
    setCredentials((prev) => ({ ...prev, [key]: value }));
  };

  const canUseVault = !!instance?.vault_secret_path;

  if (loading) {
    return (
      <div className="p-6">
        <div className="text-gray-600">Chargement...</div>
      </div>
    );
  }

  if (loadError) {
    return (
      <div className="p-6">
        <div className="flex items-center gap-2 text-red-700 bg-red-50 border border-red-100 p-3 rounded-lg">
          <XCircleIcon className="w-5 h-5" />
          <div>{loadError}</div>
        </div>
        <div className="mt-4">
          <button
            onClick={() => router.back()}
            className="btn-secondary flex items-center gap-2"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            Retour
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <button
          onClick={() => router.push(`/dashboard/aggregators/${providerId}`)}
          className="btn-secondary flex items-center gap-2"
        >
          <ArrowLeftIcon className="w-4 h-4" />
          Retour
        </button>

        <div className="text-right">
          <div className="text-sm text-gray-500">{providerDetails?.name}</div>
          <div className="text-lg font-semibold text-gray-900">
            {instance?.name}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <KeyIcon className="w-5 h-5 text-blue-600" />
            <div>
              <div className="font-semibold text-gray-900">API Credentials</div>
              <div className="text-sm text-gray-500">
                Provider: <span className="font-medium">{providerCode}</span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              disabled={savingCredentials || loadingCredentials}
              onClick={saveCredentials}
              className="btn-primary flex items-center gap-2"
            >
              <CheckCircleIcon className="w-4 h-4" />
              {savingCredentials ? "Sauvegarde..." : "Sauvegarder"}
            </button>
          </div>
        </div>

        <div className="mt-4 flex items-center justify-between gap-4 p-3 rounded-lg border border-gray-100 bg-gray-50">
          <div className="text-sm text-gray-700">
            Source: <span className="font-medium">{source}</span>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={async () => {
                setSource("db");
                await loadCredentials("db");
              }}
              className={`px-3 py-1.5 rounded-md text-sm font-medium border ${
                source === "db"
                  ? "bg-white border-blue-200 text-blue-700"
                  : "bg-transparent border-gray-200 text-gray-600"
              }`}
            >
              DB
            </button>
            <button
              disabled={!canUseVault}
              onClick={async () => {
                setSource("vault");
                await loadCredentials("vault");
              }}
              className={`px-3 py-1.5 rounded-md text-sm font-medium border ${
                source === "vault"
                  ? "bg-white border-blue-200 text-blue-700"
                  : "bg-transparent border-gray-200 text-gray-600"
              } ${!canUseVault ? "opacity-50 cursor-not-allowed" : ""}`}
              title={!canUseVault ? "vault_secret_path non configuré" : ""}
            >
              Vault
            </button>
          </div>
        </div>

        {loadingCredentials ? (
          <div className="mt-4 text-gray-600">Chargement des credentials...</div>
        ) : (
          <div className="mt-5 grid grid-cols-1 md:grid-cols-2 gap-4">
            {fields.map((f) => {
              const v = credentials[f.key] || "";
              const isSecret = !!f.secret;
              const reveal = !!showSecret[f.key];

              return (
                <div key={f.key} className="space-y-1">
                  <label className="text-sm font-medium text-gray-700">
                    {f.label}
                    {f.required ? " *" : ""}
                  </label>

                  <div className="relative">
                    <input
                      type={isSecret && !reveal ? "password" : "text"}
                      value={v}
                      onChange={(e) => onChangeField(f.key, e.target.value)}
                      className="w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-100 focus:border-blue-300 text-sm"
                      placeholder={f.required ? "Requis" : "Optionnel"}
                    />

                    {isSecret && (
                      <button
                        type="button"
                        onClick={() =>
                          setShowSecret((prev) => ({
                            ...prev,
                            [f.key]: !prev[f.key],
                          }))
                        }
                        className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-gray-500 hover:text-gray-700"
                        title={reveal ? "Masquer" : "Afficher"}
                      >
                        {reveal ? (
                          <EyeSlashIcon className="w-4 h-4" />
                        ) : (
                          <EyeIcon className="w-4 h-4" />
                        )}
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
        <div className="font-semibold text-gray-900">Instance</div>
        <div className="mt-3 grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
          <div className="p-3 rounded-lg border border-gray-100 bg-gray-50">
            <div className="text-gray-500">ID</div>
            <div className="font-mono text-gray-800 break-all">{instance?.id}</div>
          </div>
          <div className="p-3 rounded-lg border border-gray-100 bg-gray-50">
            <div className="text-gray-500">Vault Secret Path</div>
            <div className="font-mono text-gray-800 break-all">
              {instance?.vault_secret_path || "(non configuré)"}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
