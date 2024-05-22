import plotly.express as px
import pandas as pd

title = "Handshake Performance Comparison"
x_axis_label = "# of parallel clients"
y_axis_label = "Handshakes per second"
legend_label = "Handshake type"

clients = [1, 2, 4, 6, 8, 10]

df = pd.read_json("data.json").set_axis(clients)
fig = px.line(df, title=title, markers=True)
fig.update_layout(
  xaxis_title = x_axis_label,
  yaxis_title = y_axis_label,
  legend_title = legend_label
)
fig.show()