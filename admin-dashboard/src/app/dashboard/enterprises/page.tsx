"use client";

import { useState, useEffect } from "react";
import {
    BuildingOfficeIcon,
    UsersIcon,
    BanknotesIcon,
    PlusIcon
} from "@heroicons/react/24/outline";

export default function EnterprisesPage() {
    const [enterprises, setEnterprises] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Mock fetch
        setTimeout(() => {
            setEnterprises([
                { id: "1", name: "Acme Corp", type: "Service", status: "ACTIVE", employees: 12 },
                { id: "2", name: "City Transport", type: "Transport", status: "ACTIVE", employees: 450 },
                { id: "3", name: "Elite School", type: "School", status: "PENDING", employees: 30 },
            ]);
            setLoading(false);
        }, 1000);
    }, []);

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Enterprises Management</h1>
                <button className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
                    <PlusIcon className="w-5 h-5" />
                    <span>Add Enterprise</span>
                </button>
            </div>

            {/* Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-blue-100 dark:bg-blue-900 text-blue-600 rounded-full">
                            <BuildingOfficeIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Total Enterprises</p>
                            <p className="text-2xl font-bold dark:text-white">3</p>
                        </div>
                    </div>
                </div>
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-green-100 dark:bg-green-900 text-green-600 rounded-full">
                            <BanknotesIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Active Invoices</p>
                            <p className="text-2xl font-bold dark:text-white">1,250</p>
                        </div>
                    </div>
                </div>
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-purple-100 dark:bg-purple-900 text-purple-600 rounded-full">
                            <UsersIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Total Employees</p>
                            <p className="text-2xl font-bold dark:text-white">492</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead className="bg-gray-50 dark:bg-gray-700 text-xs uppercase text-gray-500 dark:text-gray-400">
                            <tr>
                                <th className="px-6 py-4 font-semibold">Name</th>
                                <th className="px-6 py-4 font-semibold">Type</th>
                                <th className="px-6 py-4 font-semibold">Status</th>
                                <th className="px-6 py-4 font-semibold">Employees</th>
                                <th className="px-6 py-4 font-semibold">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {enterprises.map((ent) => (
                                <tr key={ent.id} className="hover:bg-gray-50 dark:hover:bg-gray-750 transition">
                                    <td className="px-6 py-4 font-medium text-gray-900 dark:text-white">{ent.name}</td>
                                    <td className="px-6 py-4 text-gray-500 dark:text-gray-400">
                                        <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded text-xs">{ent.type}</span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${ent.status === 'ACTIVE' ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300' :
                                                'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300'
                                            }`}>
                                            {ent.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-gray-500 dark:text-gray-400">{ent.employees}</td>
                                    <td className="px-6 py-4">
                                        <button className="text-blue-600 hover:text-blue-800 text-sm font-medium">Manage</button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
