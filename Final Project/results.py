import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from scipy import optimize as opt

text = "\
250&    0.049 &   0.007 &  0.022&\
500&    0.395 &   0.034 &  0.072&\
750&    1.354 &   0.105 &  0.164&\
1000&    2.920 &   0.217 &  0.300&\
1250&    6.118 &   0.435 &  0.501&\
1500&   10.103 &   0.723 &   0.786&\
1750&   34.168 &   2.177 &   1.439&\
2000&   41.578 &   2.874 &   1.902&\
2250&   77.667 &   4.555 &   3.137&\
2500&   93.060 &   5.463 &   3.694&\
5000&  820.133 &  43.913 &  25.819&\
7500& 3573.429 & 175.860 &  95.659&\
10000& 6558.381 & 353.265 & 203.442&\
15000&     *    &1318.689 & 702.739&\
20000&     *    &9400.968 & 3656.001\
"

parsed = text.split('&')
print(parsed)
results = []
for i, t in enumerate(parsed):
	if '*' in t:
		results.append(0)
	else:
		results.append(float(t))

print(results)
n, row, go, distributed = [], [], [], []
for i, _ in enumerate(results):
	if   i % 4 == 0:
		n.append(_)
	elif i % 4 == 1:
		row.append(_)
	elif i % 4 == 2:
		go.append(_)
	elif i % 4 == 3:
		distributed.append(_)

def objective(x, a, b):
	return a*x + b

def exponential(x, a, b):
	return a*x**b

def fit(x_values, y_values):
	popt, _ = opt.curve_fit(exponential, x_values, y_values)
	a, b = popt
	x_new = np.linspace(min(x_values), max(x_values), 10)
	y_new = exponential(x_new, a, b)
	return x_new, y_new


plt.figure(figsize=(10,6))
plt.loglog(n[:-2], row[:-2], 'ro', label='Row-wise')
# n_new, row_hat = fit(n, row)
# plt.plot(n_new, row_hat, 'r--', label='Row-wise Fit')
plt.loglog(n, go, 'gx', label='Goroutines')
# _, go_hat = fit(n, go)
# plt.plot(n_new, go_hat, 'g--', label='Goroutines Fit')
plt.loglog(n, distributed, 'bv', label='Distributed')
# _, distributed_hat = fit(n, distributed)
# plt.plot(n_new, distributed_hat, 'b--', label='Distributed Fit')

plt.legend(shadow=True)
plt.grid(alpha=0.4618)
plt.title("Log-log plot of the timings vs. dimension")
plt.xlabel('$\\log(n)$')
plt.ylabel('Time ($\\log$ seconds)')
plt.savefig('./imgs/results.png', bbox_inches="tight", dpi=300)
plt.show()