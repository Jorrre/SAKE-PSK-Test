import plotly.express as px
import pandas as pd

title = "Handshake Performance Comparison"
x_axis_title = "# of parallel clients"
y_axis_title = "Handshakes per second"
legend_title = "Handshake type"

clients = [1, 2, 4, 6, 8, 10]

loopback_df = pd.read_json("loopback.json").set_axis(clients)
loopback_fig = px.line(loopback_df, title="Loopback", markers=True)
loopback_fig.update_layout(
  xaxis_title = x_axis_title,
  yaxis_title = y_axis_title,
  legend_title = legend_title
)
loopback_fig.show()
  
vpn_df = pd.read_json("vpn.json").set_axis(clients)
vpn_fig = px.line(vpn_df, title="VPN", markers=True)
vpn_fig.update_layout(
  xaxis_title = x_axis_title,
  yaxis_title = y_axis_title,
  legend_title = legend_title
)
vpn_fig.show()