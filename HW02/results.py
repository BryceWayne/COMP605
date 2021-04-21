import pandas as pd
import matplotlib.pyplot as plt

results = [{'n':  250, 'row':  0.044, 'col':   0.040, 'go': 0.005},
           {'n':  500, 'row':  0.359, 'col':   0.341, 'go': 0.031},
           {'n':  750, 'row':  1.381, 'col':   2.407, 'go': 0.140},
           {'n': 1000, 'row':  3.823, 'col':  10.804, 'go': 0.490},
           {'n': 1250, 'row': 10.740, 'col':  25.997, 'go': 1.161},
           {'n': 1500, 'row': 20.173, 'col':  47.134, 'go': 2.224},
           {'n': 1750, 'row': 45.961, 'col':  77.037, 'go': 4.136},
           {'n': 2000, 'row': 55.245, 'col': 114.245, 'go': 5.392},]

df = pd.DataFrame(results)
plt.figure(figsize=(10,6))
plt.plot(df['n'].to_numpy(), df['row'].to_numpy(), c='#27ff27', label='Row')
plt.plot(df['n'].to_numpy(), df['col'].to_numpy(), c='#2727ff', label='Col')
plt.plot(df['n'].to_numpy(), df['go'].to_numpy(), c='#ff2727', label='Goroutine')
plt.xlabel('$n$')
plt.ylabel('seconds')
plt.xlim(250,2000)
plt.ylim(0, 135)
plt.legend(shadow=True)
plt.grid(alpha=0.4618)
plt.title('Golang Matrix-Matrix Multiplication')
plt.savefig('timingsGo.png', bbox_inches='tight', dpi=300)
plt.show()