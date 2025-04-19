go test -bench=. -benchmem > benchmark_results_baseline.txt

go test -bench=. -benchmem > benchmark_results_optimized_s1.txt

benchstat benchmark_results_baseline.txt benchmark_results_optimized_s1.txt > benchmark_compare_baseline_vs_strategy_1.txt