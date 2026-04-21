"use client";

import { useState, useEffect } from "react";
import { 
  Plus, 
  Search, 
  Filter, 
  RefreshCw 
} from "lucide-react";

const departments = ["ICU", "Emergency", "Surgery", "Pediatrics", "Oncology"];
const categories = ["Ventilators", "Vital Monitors", "Infusion Pumps", "Defibrillators"];

export default function AvailabilityBoard() {
  const [isRefreshing, setIsRefreshing] = useState(false);

  // Mock real-time update effect
  useEffect(() => {
    const interval = setInterval(() => {
      // Logic for real-time flashes/updates would go here
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-8 animate-fade-in">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold text-white">Live Availability</h2>
          <p className="text-slate-400 mt-2">Real-time equipment status across all departments.</p>
        </div>
        <div className="flex space-x-4">
          <div className="glass flex items-center px-4 py-2 rounded-xl border border-white/5 focus-within:border-indigo-500/50 transition-all">
            <Search size={18} className="text-slate-400 mr-3" />
            <input 
              type="text" 
              placeholder="Search category..." 
              className="bg-transparent border-none outline-none text-sm w-64"
            />
          </div>
          <button 
            onClick={() => {
              setIsRefreshing(true);
              setTimeout(() => setIsRefreshing(false), 1000);
            }}
            className="btn-secondary flex items-center space-x-2"
          >
            <RefreshCw size={18} className={isRefreshing ? "animate-spin" : ""} />
            <span>Sync</span>
          </button>
        </div>
      </header>

      {/* Grid Container */}
      <div className="glass rounded-3xl overflow-hidden border border-white/5">
        <div className="overflow-x-auto">
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-white/5">
                <th className="p-6 text-left text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-white/5">
                  Category
                </th>
                {departments.map(dept => (
                  <th key={dept} className="p-6 text-center text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-white/5">
                    {dept}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {categories.map((cat, i) => (
                <tr key={cat} className="hover:bg-white/[0.02] transition-colors group">
                  <td className="p-6 font-semibold text-white border-r border-white/5">
                    <div className="flex items-center space-x-3">
                      <div className="w-2 h-2 rounded-full bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.5)]" />
                      <span>{cat}</span>
                    </div>
                  </td>
                  {departments.map((dept, j) => {
                    // Mock data logic
                    const total = Math.floor(Math.random() * 10) + 2;
                    const available = Math.floor(Math.random() * (total + 1));
                    const isLow = available <= 1;

                    return (
                      <td key={`${cat}-${dept}`} className="p-6 text-center">
                        <div className={`inline-flex flex-col items-center justify-center w-16 h-16 rounded-2xl transition-all duration-300 ${
                          isLow 
                            ? "bg-rose-500/10 border border-rose-500/20 text-rose-500 shadow-lg shadow-rose-500/10" 
                            : "bg-emerald-500/5 border border-white/5 text-emerald-500 group-hover:border-white/10"
                        }`}>
                          <span className="text-xl font-bold">{available}</span>
                          <span className="text-[10px] uppercase font-bold opacity-60">/ {total}</span>
                        </div>
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="flex justify-end space-x-6">
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 rounded bg-emerald-500/20 border border-emerald-500/40" />
          <span className="text-xs text-slate-400">High Availability</span>
        </div>
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 rounded bg-rose-500/20 border border-rose-500/40" />
          <span className="text-xs text-slate-400">Low/Critical Stock</span>
        </div>
      </div>
    </div>
  );
}
