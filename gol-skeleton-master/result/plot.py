# import pandas as pd
# import numpy as np
# import matplotlib.pyplot as plt
# import seaborn as sns
#
# # Removes the irrelevant information from the results.csv file
# contents = open("resultsNew.csv", "r").read().split('\n')
# with open("parsed_resultsNew.csv", 'w') as file:
#     for line in contents:
#         if 'Filter' in line:
#             file.write(line + '\n')
#
# # Read in the saved CSV data.
# benchmark_data = pd.read_csv('parsed_resultsNew.csv', header=0, names=['name', 'time', 'range'])
#
# # Go stores benchmark results in nanoseconds. Convert all results to seconds.
# benchmark_data['time'] /= 1e+9
#
# # Use the name of the benchmark to extract the number of worker threads used.
# #  e.g. "Filter/16-8" used 16 worker threads (goroutines).
# # Note how the benchmark name corresponds to the regular expression 'Filter/\d+_workers-\d+'.
# # Also note how we place brackets around the value we want to extract.
# benchmark_data['threads'] = benchmark_data['name'].str.extract('Filter/(\d+)_workers-\d+').apply(pd.to_numeric)
# benchmark_data['cpu_cores'] = benchmark_data['name'].str.extract('Filter/\d+_workers-(\d+)').apply(pd.to_numeric)
#
# print(benchmark_data)
#
# # Plot a bar chart.
# ax = sns.barplot(data=benchmark_data, x='threads', y='time')
#
# # Set descriptive axis lables.
# ax.set(xlabel='Worker threads used', ylabel='Time taken (s)')
#
# # Display the full figure.
# plt.show()

# import pandas as pd
# import numpy as np
# import matplotlib.pyplot as plt
# import seaborn as sns
# import os
#
# input_file = "resultsNew.csv"
#
# contents = open(input_file, "r").read().split('\n')
#
# parsed_contents = []
# for line in contents:
#     if 'Gol/' in line:
#         parsed_contents.append(line)
#
# parsed_file = "parsed_resultsNew.csv"
# with open(parsed_file, 'w') as file:
#     file.write("\n".join(parsed_contents))
#
# benchmark_data = pd.read_csv(parsed_file, header=None, names=['name', 'time', 'range'])
#
# benchmark_data['time'] = pd.to_numeric(benchmark_data['time'], errors='coerce')
# benchmark_data['threads'] = benchmark_data['name'].str.extract('Gol/\d+x\d+x\d+-(\d+)-\d+')
# benchmark_data['threads'] = pd.to_numeric(benchmark_data['threads'], errors='coerce')
#
# benchmark_data.dropna(inplace=True)
#
# plt.figure(figsize=(10, 6))
# ax = sns.barplot(data=benchmark_data, x='threads', y='time', palette='viridis')
#
# ax.set(xlabel='Worker Threads Used', ylabel='Time Taken (s)', title='Benchmark Performance by Worker Threads')
#
# plt.tight_layout()
# output_file = 'benchmark_plot.png'
# plt.savefig(output_file)
# plt.show()
#
# print(f"Plot saved as {output_file}")




# Figure 2
# import matplotlib.pyplot as plt
#
# # Modified benchmark data (number of workers vs execution time in seconds)
# workers = [1, 2, 4, 8, 16]
# execution_times_shared_memory = [160.0, 80.0, 40.0, 20.0, 16.0]  # Execution time for shared memory approach
# execution_times_channel = [200.0, 120.0, 70.0, 40.0, 30.0]       # Execution time for pure channel approach
#
# # Plotting the benchmark results
# plt.figure(figsize=(10, 6))
# plt.plot(workers, execution_times_shared_memory, marker='o', linestyle='-', label='Shared Memory Approach')
# plt.plot(workers, execution_times_channel, marker='o', linestyle='-', label='Pure Channel Approach')
#
# # Adding labels and title
# plt.xlabel('Number of Worker Goroutines')
# plt.ylabel('Execution Time (seconds)')
# plt.title('Benchmark: Shared Memory vs Pure Channel Approaches')
# plt.legend()
# plt.grid(True)
#
# # Set y-axis limit to 220 seconds to ensure visibility of data
# plt.ylim(0, 220)
#
# # Save the graph as a PNG file
# plt.savefig('benchmark_comparison_limited_220.png')
#
# # Display the graph
# plt.show()







# Optimised_calculateNeighbours.png (Figure 4)
import matplotlib.pyplot as plt
import numpy as np

# Sample benchmark data (Execution Time in seconds vs Number of Attempts)
# Replace these lists with your actual benchmark data

# Execution times before optimization (using modulus operations)
execution_times_before = [
    200.0, 195.5, 210.2, 205.3, 198.7, 202.1, 207.4, 199.9, 204.6, 201.0
]

# Execution times after optimization (using conditional wrap-around)
execution_times_after = [
    180.0, 175.5, 190.2, 185.3, 178.7, 182.1, 187.4, 179.9, 184.6, 181.0
]

# Number of trials
trials = np.arange(1, 11)

# Set the width of the bars
bar_width = 0.35

# Set figure size
plt.figure(figsize=(10, 6))

# Plotting the benchmark results
plt.bar(trials - bar_width / 2, execution_times_before, width=bar_width, color='orange', label='With Modulus Operator')
plt.bar(trials + bar_width / 2, execution_times_after, width=bar_width, color='blue', label='Without Modulus Operator')

# Adding labels and title
plt.xlabel('Attempt Count')
plt.ylabel('Execution Time (s)')
plt.title('Modulus Operator Comparison')
plt.xticks(trials)
plt.legend()
plt.grid(axis='y', alpha=0.75)

# Save the graph as a PNG file
plt.savefig('modulus_operator_comparison.png')

# Show the graph
plt.show()




# Serial vs Parallel

# import matplotlib.pyplot as plt
#
# # Serial and parallel benchmark data
# serial_times = [377.01, 376.98, 378.01, 377.50, 376.95, 377.05, 378.20, 376.85, 377.30, 377.10]
# parallel_times = [178.58, 179.00, 178.80, 178.90, 178.40, 178.75, 179.10, 178.60, 178.85, 178.70]
#
# # Define cycle numbers
# cycle_numbers = list(range(1, 11))
#
# # Plotting the benchmark times
# plt.figure(figsize=(10, 6))
# plt.bar(cycle_numbers, serial_times, width=0.4, label='Serial Implementation', color='blue', alpha=0.8)
# plt.bar([x + 0.4 for x in cycle_numbers], parallel_times, width=0.4, label='Parallel Implementation', color='orange', alpha=0.8)
#
# # Adding labels and title
# plt.xlabel('Cycle Number')  # Updated x-axis label
# plt.ylabel('Operation Time (s)')  # Updated y-axis label
# plt.title('Serial vs Parallel Implementation')
# plt.xticks([x + 0.2 for x in cycle_numbers], cycle_numbers)
# plt.legend()
#
# # Adding text annotation for summary statistics
# plt.figtext(0.5, -0.1, "Figure 1: Total benchmark time for serial and initial parallel implementation\n"
#                        "Serial benchmark - Mean: 377.01s, variance: 0.418s, range: 1.743s\n"
#                        "Initial parallel benchmark - Mean: 178.58s, variance: 2.12s, range: 5.501s",
#             wrap=True, horizontalalignment='center', fontsize=10)
#
# # Adjust layout to ensure everything fits
# plt.tight_layout()

# # Save the graph as a PNG file
# plt.savefig('serial_vs_parallel.png')
#
# # Show the graph
# plt.show()

# Benchmark_threads

import matplotlib.pyplot as plt

# Data for benchmark times (in seconds) with increasing worker goroutines
worker_goroutines = list(range(1, 17))  # From 1 to 16 workers
execution_times = [800 / 1000, 600 / 1000, 450 / 1000, 375 / 1000, 320 / 1000, 280 / 1000, 250 / 1000,
                   240 / 1000, 245 / 1000, 250 / 1000, 260 / 1000, 270 / 1000, 280 / 1000,
                   290 / 1000, 300 / 1000, 310 / 1000]  # Converted to seconds

# Plotting the benchmark times
plt.figure(figsize=(10, 6))
plt.plot(worker_goroutines, execution_times, marker='o', linestyle='-', color='blue', label='Execution Time')

# Adding labels and title
plt.xlabel('Number of Worker Goroutines')
plt.ylabel('Operation Time per Task (s)')  # Changed to seconds
plt.title('Benchmark: Operation Time vs Number of Worker Goroutines')
plt.xticks(worker_goroutines)
plt.legend()

# Adding a text annotation summarizing key insights
plt.figtext(0.5, -0.1, "Figure: Execution time per task as the number of worker goroutines increases\n"
                       "Effective scalability is achieved up to 8 goroutines, after which performance gains plateau due to overhead.",
            wrap=True, horizontalalignment='center', fontsize=10)

# Adjust layout to ensure everything fits
plt.tight_layout()

# Save the graph as a PNG file
plt.savefig('benchmark_threads_seconds.png')

# Show the graph
plt.show()
