# Global Resouce Service logging guidelines and conventions
We initiate the service logging guidelines and conventions early in the project stage with manageable code base for now. This will be a starting point for GRS developers to follow when they join the team or contribute to the project. all of those guidelines can be refined later stages in the project.

**Overall the goal for logging are:**
1. provide key code flow to whom reads the log for either looking into info or debugging
2. provide leveled information for "essential logging info", "verbose logging info" and "debug logging info"

**Initial set of guidelines:**
1. log level in GRS has roughly 3 levels, V(4), V(6) and V(9), in addition to the essential logging level ( where no V() is added ).
2. errors, warning in code should be logged out with Errorf() or Warningf(). some error info can be logged with Infof() if the errors is interpreted in the code with info.
3. V(4) logging provides key info to indicate KEY code flow such as interface calls, and to indicate the service running status, essential summary of perf info. this should be kept at minimal to avoid log overflow too much in production env. most functional bug can be investigated or provide further hints to investigate the functional issues.
4. V(6) logging provides so most or all code flows for helper or internal functions etc. some level of perf info or info to investigate functional or perf issues can be added with V(6). however, large volume of log entry, or logs in loops should be well-considered to only log needed info in the data structures.
5. V(9) or above is considered as debugging log level, where detailed traces can be added at this level.
6. essential logging info will just use Klog.Infof(), i.e. for 2 or below, just use infof()
7. most test and recommended production log level is V4.
8. V7, V8 or 9+ are not used for now.