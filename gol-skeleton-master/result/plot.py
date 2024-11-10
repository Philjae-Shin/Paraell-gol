import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

# Removes the irrelevant information from the results.csv file
contents = open("resultsNew.csv", "r").read().split('\n')
with open("parsed_resultsNew.csv", 'w') as file:
    for line in contents:
        if 'Filter' in line:
            file.write(line + '\n')

# Read in the saved CSV data.
benchmark_data = pd.read_csv('parsed_resultsNew.csv', header=0, names=['name', 'time', 'range'])

# Go stores benchmark results in nanoseconds. Convert all results to seconds.
benchmark_data['time'] /= 1e+9

# Use the name of the benchmark to extract the number of worker threads used.
#  e.g. "Filter/16-8" used 16 worker threads (goroutines).
# Note how the benchmark name corresponds to the regular expression 'Filter/\d+_workers-\d+'.
# Also note how we place brackets around the value we want to extract.
benchmark_data['threads'] = benchmark_data['name'].str.extract('Filter/(\d+)_workers-\d+').apply(pd.to_numeric)
benchmark_data['cpu_cores'] = benchmark_data['name'].str.extract('Filter/\d+_workers-(\d+)').apply(pd.to_numeric)

print(benchmark_data)

# Plot a bar chart.
ax = sns.barplot(data=benchmark_data, x='threads', y='time')

# Set descriptive axis lables.
ax.set(xlabel='Worker threads used', ylabel='Time taken (s)')

# Display the full figure.
plt.show()

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import os

input_file = "resultsNew.csv"

contents = open(input_file, "r").read().split('\n')

parsed_contents = []
for line in contents:
    if 'Gol/' in line:
        parsed_contents.append(line)

parsed_file = "parsed_resultsNew.csv"
with open(parsed_file, 'w') as file:
    file.write("\n".join(parsed_contents))

benchmark_data = pd.read_csv(parsed_file, header=None, names=['name', 'time', 'range'])

benchmark_data['time'] = pd.to_numeric(benchmark_data['time'], errors='coerce')
benchmark_data['threads'] = benchmark_data['name'].str.extract('Gol/\d+x\d+x\d+-(\d+)-\d+')
benchmark_data['threads'] = pd.to_numeric(benchmark_data['threads'], errors='coerce')

benchmark_data.dropna(inplace=True)

plt.figure(figsize=(10, 6))
ax = sns.barplot(data=benchmark_data, x='threads', y='time', palette='viridis')

ax.set(xlabel='Worker Threads Used', ylabel='Time Taken (s)', title='Benchmark Performance by Worker Threads')

plt.tight_layout()
output_file = 'benchmark_plot.png'
plt.savefig(output_file)
plt.show()

print(f"Plot saved as {output_file}")



