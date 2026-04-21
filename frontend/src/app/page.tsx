import { 
  TrendingUp, 
  AlertCircle, 
  History 
} from "lucide-react";

const stats = [
  { label: "Overall Utilization", value: "78%", change: "+12%", trend: "up", color: "indigo" },
  { label: "In Maintenance", value: "14", change: "-2", trend: "down", color: "rose" },
  { label: "Sharing Efficiency", value: "92%", change: "+5%", trend: "up", color: "emerald" },
  { label: "Pending Requests", value: "8", change: "New", trend: "up", color: "amber" },
];

export default function Dashboard() {
  return (
    <div className="space-y-8 animate-fade-in">
      <header>
        <h2 className="text-3xl font-bold text-white">Hospital Overview</h2>
        <p className="text-slate-400 mt-2">Real-time utilization metrics and system health.</p>
      </header>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => (
          <div key={stat.label} className="glass-card">
            <div className="flex justify-between items-start">
              <span className="text-slate-400 text-sm font-medium">{stat.label}</span>
              <div className={`p-2 rounded-lg bg-${stat.color}-500/10 text-${stat.color}-500`}>
                {stat.trend === "up" ? <TrendingUp size={16} /> : <AlertCircle size={16} />}
              </div>
            </div>
            <div className="mt-4 flex items-end justify-between">
              <h3 className="text-2xl font-bold text-white">{stat.value}</h3>
              <span className={`text-xs font-semibold ${stat.trend === "up" ? "text-emerald-500" : "text-rose-500"}`}>
                {stat.change}
              </span>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Chart Area */}
        <div className="lg:col-span-2 glass-card h-[400px] flex flex-col">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Utilization Trends</h3>
            <button className="text-indigo-400 text-sm hover:underline">View Full Report</button>
          </div>
          <div className="flex-1 flex items-end space-x-4 px-2 pb-2">
            {[40, 65, 45, 80, 55, 90, 75, 85, 60, 95, 70, 80].map((h, i) => (
              <div key={i} className="flex-1 group relative">
                <div 
                  className="bg-indigo-500/20 group-hover:bg-indigo-500/40 transition-all rounded-t-lg"
                  style={{ height: `${h}%` }}
                ></div>
                <div className="absolute -top-8 left-1/2 -translate-x-1/2 bg-white text-slate-900 text-[10px] px-1.5 py-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity">
                  {h}%
                </div>
              </div>
            ))}
          </div>
          <div className="flex justify-between mt-4 text-[10px] text-slate-500 px-1">
            <span>JAN</span><span>MAR</span><span>MAY</span><span>JUL</span><span>SEP</span><span>NOV</span>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="glass-card">
          <div className="flex items-center space-x-2 mb-6">
            <History size={18} className="text-indigo-400" />
            <h3 className="text-lg font-bold">Recent Activity</h3>
          </div>
          <div className="space-y-6">
            {[
              { type: "Request", desc: "ICU requested Vent-04 from Surgery", time: "2m ago" },
              { type: "Maintenance", desc: "Monitor-12 flagged for calibration", time: "15m ago" },
              { type: "Handoff", desc: "Vent-04 handoff confirmed by ER", time: "45m ago" },
              { type: "System", desc: "Daily utilization report generated", time: "1h ago" },
            ].map((act, i) => (
              <div key={i} className="flex flex-col space-y-1">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-white">{act.type}</span>
                  <span className="text-slate-500 text-xs">{act.time}</span>
                </div>
                <p className="text-xs text-slate-400 leading-relaxed">{act.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
