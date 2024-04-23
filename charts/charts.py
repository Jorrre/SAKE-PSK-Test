import plotly.express as px
import pandas as pd

df = pd.DataFrame({
  '# of clients': range(1, 11),
  'Full handshake': [147.1, 281.6, 383.7, 460.7, 525.0, 588.6, 645.9, 696.7, 768.3, 822.6],
  'Plain PSK': [3791.3, 7992.0, 13279.6, 14412.8, 12854.9, 6277.7, 5345.4, 5172.4, 5185.1, 5263.1],
  'ECDHE-PSK': [1741.2, 4095.7, 6604.1, 7940.0, 8767.5, 8142.4, 4102.2, 3840.3, 4205.4, 7241.2],
  'SAKE-PSK': [3634.7, 7386.7, 11850.4, 12274.6, 9295.4, 5633.8, 5700.4, 5766.9, 4864.9, 5438.1],
})

df_melted = df.melt(id_vars='# of clients', var_name='Line', value_name='Handshakes per second')

fig = px.line(df_melted, x='# of clients', y='Handshakes per second', color='Line', title='Full handshake vs PSK handshake comparison', markers=True)

fig.show()