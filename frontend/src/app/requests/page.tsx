"use client";

import { 
  Plus, 
  ArrowRight, 
  CheckCircle2, 
  XCircle, 
  Clock,
  MoreVertical
} from "lucide-react";

const requests = [
  { id: "1", category: "Ventilator", from: "ICU", to: "Emergency", status: "in_transit", urgency: "high", time: "12m ago" },
  { id: "2", category: "Infusion Pump", from: "Pediatrics", to: "Surgery", status: "pending", urgency: "normal", time: "45m ago" },
  { id: "3", category: "Monitor", from: "Oncology", to: "ICU", status: "active", urgency: "emergency", time: "1h ago" },
  { id: "4", category: "Defibrillator", from: "Surgery", to: "ER", status: "return_pending", urgency: "normal", time: "3h ago" },
];

export default function SharingRequests() {
  return (
    <div className="space-y-8 animate-fade-in">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold text-white">Sharing Workflow</h2>
          <p className="text-slate-400 mt-2">Manage inter-departmental equipment sharing requests.</p>
        </div>
        <button className="btn-primary flex items-center space-x-2">
          <Plus size={18} />
          <span>New Request</span>
        </button>
      </header>

      <div className="grid grid-cols-1 gap-4">
        {requests.map((req) => (
          <div key={req.id} className="glass-card hover:translate-x-1">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-6">
                {/* Urgency Indicator */}
                <div className={`w-1 h-12 rounded-full ${
                  req.urgency === "emergency" ? "bg-rose-500 shadow-[0_0_12px_rgba(244,63,94,0.5)]" :
                  req.urgency === "high" ? "bg-amber-500 shadow-[0_0_12px_rgba(245,158,11,0.5)]" :
                  "bg-indigo-500 shadow-[0_0_12px_rgba(99,102,241,0.5)]"
                }`} />

                <div>
                  <h4 className="text-lg font-bold text-white">{req.category}</h4>
                  <div className="flex items-center space-x-2 mt-1">
                    <span className="text-slate-400 text-sm">{req.from}</span>
                    <ArrowRight size={14} className="text-slate-600" />
                    <span className="text-slate-400 text-sm">{req.to}</span>
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-8">
                <div className="flex flex-col items-end">
                  <span className={`status-badge ${
                    req.status === "active" ? "status-available" :
                    req.status === "in_transit" ? "status-in-use" :
                    "status-maintenance"
                  }`}>
                    {req.status.replace("_", " ")}
                  </span>
                  <span className="text-[10px] text-slate-500 mt-1 uppercase font-bold tracking-wider">
                    {req.time}
                  </span>
                </div>

                <div className="flex items-center space-x-2">
                  {req.status === "pending" ? (
                    <>
                      <button className="p-2 text-rose-400 hover:bg-rose-400/10 rounded-lg transition-colors">
                        <XCircle size={20} />
                      </button>
                      <button className="p-2 text-emerald-400 hover:bg-emerald-400/10 rounded-lg transition-colors">
                        <CheckCircle2 size={20} />
                      </button>
                    </>
                  ) : (
                    <button className="p-2 text-slate-400 hover:bg-white/5 rounded-lg transition-colors">
                      <MoreVertical size={20} />
                    </button>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Summary Footer */}
      <div className="glass p-6 rounded-2xl flex justify-between items-center border border-white/5">
        <div className="flex space-x-12">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-indigo-500/10 text-indigo-400">
              <Clock size={18} />
            </div>
            <div>
              <p className="text-[10px] text-slate-500 uppercase font-bold">Avg. Match Time</p>
              <p className="text-lg font-bold">14 mins</p>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-emerald-500/10 text-emerald-400">
              <CheckCircle2 size={18} />
            </div>
            <div>
              <p className="text-[10px] text-slate-500 uppercase font-bold">Success Rate</p>
              <p className="text-lg font-bold">98.4%</p>
            </div>
          </div>
        </div>
        <button className="text-indigo-400 text-sm font-semibold hover:underline">
          View History Archive
        </button>
      </div>
    </div>
  );
}
