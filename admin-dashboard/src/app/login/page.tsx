'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { login } from '@/lib/api';
import { useAuthStore } from '@/stores/authStore';
import {
    LockClosedIcon,
    EnvelopeIcon,
    EyeIcon,
    EyeSlashIcon,
    ShieldCheckIcon,
    ArrowRightIcon,
    SparklesIcon,
} from '@heroicons/react/24/outline';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const setAdmin = useAuthStore((state) => state.setAdmin);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const data = await login(email, password);
            setAdmin(data.admin, data.admin.role?.permissions?.map((p: any) => p.code) || []);
            router.push('/dashboard');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Identifiants incorrects. Veuillez réessayer.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex">
            {/* Left side - Decorative */}
            <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden bg-gradient-to-br from-slate-900 via-indigo-900 to-purple-900">
                {/* Animated Background Elements */}
                <div className="absolute inset-0">
                    <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-indigo-600/30 rounded-full blur-3xl animate-pulse" />
                    <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-purple-600/30 rounded-full blur-3xl animate-pulse delay-1000" />
                    <div className="absolute top-1/2 left-1/2 w-64 h-64 bg-pink-600/20 rounded-full blur-3xl animate-pulse delay-500" />
                </div>

                {/* Grid Pattern Overlay */}
                <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-10" />

                {/* Content */}
                <div className="relative z-10 flex flex-col justify-center p-16 text-white">
                    <div className="max-w-md">
                        {/* Logo */}
                        <div className="flex items-center gap-4 mb-12">
                            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-2xl shadow-indigo-500/30">
                                <img src="/logo.png" alt="Zekora" className="w-10 h-10 object-contain" />
                            </div>
                            <div>
                                <h1 className="text-3xl font-bold">Zekora</h1>
                                <p className="text-indigo-300">Administration</p>
                            </div>
                        </div>

                        <h2 className="text-4xl font-bold mb-6 leading-tight">
                            Gérez votre plateforme<br />
                            <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-purple-400">
                                en toute simplicité
                            </span>
                        </h2>

                        <p className="text-lg text-slate-300 mb-10">
                            Accédez à votre tableau de bord administrateur pour superviser
                            les opérations, gérer les utilisateurs et assurer la sécurité de la plateforme.
                        </p>

                        {/* Features */}
                        <div className="space-y-4">
                            {[
                                { icon: ShieldCheckIcon, text: 'Sécurité renforcée & 2FA' },
                                { icon: SparklesIcon, text: 'Interface intuitive' },
                                { icon: LockClosedIcon, text: 'Accès contrôlé par rôles' },
                            ].map((feature, i) => (
                                <div key={i} className="flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-xl bg-white/10 flex items-center justify-center">
                                        <feature.icon className="w-5 h-5 text-indigo-400" />
                                    </div>
                                    <span className="text-slate-200">{feature.text}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>

            {/* Right side - Login Form */}
            <div className="flex-1 flex items-center justify-center p-8 bg-gradient-to-br from-slate-50 via-white to-slate-100">
                <div className="w-full max-w-md">
                    {/* Mobile Logo */}
                    <div className="lg:hidden text-center mb-10">
                        <div className="inline-flex items-center gap-3">
                            <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/25">
                                <img src="/logo.png" alt="Zekora" className="w-8 h-8 object-contain" />
                            </div>
                            <div className="text-left">
                                <h1 className="text-2xl font-bold text-slate-900">Zekora</h1>
                                <p className="text-indigo-600 text-sm font-medium">Admin Panel</p>
                            </div>
                        </div>
                    </div>

                    {/* Card */}
                    <div className="bg-white rounded-3xl shadow-2xl shadow-slate-200/50 p-8 border border-slate-100">
                        <div className="text-center mb-8">
                            <div className="hidden lg:flex items-center justify-center w-16 h-16 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 shadow-lg shadow-indigo-500/25">
                                <LockClosedIcon className="w-8 h-8 text-white" />
                            </div>
                            <h2 className="text-2xl font-bold text-slate-900">Connexion Admin</h2>
                            <p className="text-slate-500 mt-2">Entrez vos identifiants pour accéder au panneau</p>
                        </div>

                        {/* Error Alert */}
                        {error && (
                            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-3 animate-slide-up">
                                <div className="w-8 h-8 rounded-lg bg-red-100 flex items-center justify-center flex-shrink-0">
                                    <svg className="w-4 h-4 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </div>
                                <div>
                                    <p className="text-sm font-semibold text-red-800">Erreur de connexion</p>
                                    <p className="text-sm text-red-600">{error}</p>
                                </div>
                            </div>
                        )}

                        <form onSubmit={handleSubmit} className="space-y-5">
                            {/* Email Field */}
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">
                                    Adresse email
                                </label>
                                <div className="relative">
                                    <EnvelopeIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                                    <input
                                        type="email"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        className="w-full pl-12 pr-4 py-3.5 rounded-xl border-2 border-slate-200 focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none transition-all bg-slate-50 focus:bg-white"
                                        placeholder="admin@zekora.com"
                                        required
                                    />
                                </div>
                            </div>

                            {/* Password Field */}
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">
                                    Mot de passe
                                </label>
                                <div className="relative">
                                    <LockClosedIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        className="w-full pl-12 pr-12 py-3.5 rounded-xl border-2 border-slate-200 focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none transition-all bg-slate-50 focus:bg-white"
                                        placeholder="••••••••••••"
                                        required
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        className="absolute right-4 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-slate-600 transition-colors"
                                    >
                                        {showPassword ? (
                                            <EyeSlashIcon className="w-5 h-5" />
                                        ) : (
                                            <EyeIcon className="w-5 h-5" />
                                        )}
                                    </button>
                                </div>
                            </div>

                            {/* Remember & Forgot */}
                            <div className="flex items-center justify-between text-sm">
                                <label className="flex items-center gap-2 cursor-pointer group">
                                    <input type="checkbox" className="w-4 h-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500" />
                                    <span className="text-slate-600 group-hover:text-slate-800">Se souvenir de moi</span>
                                </label>
                                <a href="#" className="text-indigo-600 hover:text-indigo-700 font-medium">
                                    Mot de passe oublié ?
                                </a>
                            </div>

                            {/* Submit Button */}
                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full flex items-center justify-center gap-2 py-4 px-6 rounded-xl font-semibold text-white bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-lg shadow-indigo-500/25 hover:shadow-xl hover:shadow-indigo-500/30 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none"
                            >
                                {loading ? (
                                    <>
                                        <span className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                                        Connexion en cours...
                                    </>
                                ) : (
                                    <>
                                        Se connecter
                                        <ArrowRightIcon className="w-5 h-5" />
                                    </>
                                )}
                            </button>
                        </form>
                    </div>

                    {/* Footer */}
                    <div className="mt-8 text-center">
                        <p className="text-sm text-slate-500">
                            Accès réservé aux administrateurs autorisés
                        </p>
                        <div className="flex items-center justify-center gap-2 mt-2 text-xs text-slate-400">
                            <ShieldCheckIcon className="w-4 h-4" />
                            Connexion sécurisée SSL/TLS
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
