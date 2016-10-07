## You can Sonar!

1. Create `sonar-project.properties` file in your repository's root:
```
sonar.sources=src/main/java                # where are your source files? (relative to the repo root).
sonar.tests=src/test/java                  # Where are your test files? (relative to the repo root).
sonar.jacoco.reportPath=target/jacoco.exec # Where is the CodeCoverage report? We'll need to figure out another reporting mechanism for non-JaCoCo projects.
sonar.binaries=target/classes              # I think this only applies to java projects.
sonar.redmine.project-key=13               # I assume this is the redmine key for all piazza repos, other teams may have a different key.
sonar.sourceEncoding=UTF-8                 # You probably won't need to worry about this one.
```

2. Add a `sonar` step to your build pipeline in [`venicegeo/jenkins/Repos.groovy`](https://github.com/venicegeo/jenkins/blob/master/Repos.groovy)
```
diff --git a/Repos.groovy b/Repos.groovy
--- a/Repos.groovy
+++ b/Repos.groovy
@@ -55,7 +55,7 @@ class Repos {
     ],[
       reponame: 'pz-gateway',
       team: 'piazza',
-      pipeline: ['ionchannel_pom', 'archive', 'cf_push_int', 'cf_bg_deploy_int', 'int-release', 'run_integration_tests', 'cf_push_stage', 'cf_bg_deploy_stage', 'stage-release']
+      pipeline: ['sonar', 'ionchannel_pom', 'archive', 'cf_push_int', 'cf_bg_deploy_int', 'int-release', 'run_integration_tests', 'cf_push_stage', 'cf_bg_deploy_stage', 'stage-release']
     ],[
       reponame: 'pz-gocommon',
       lib: true,
```

3. You may need to create a `ci/sonar.sh` script in your repository; some projects may require you to generate a report outside of the sonar-runner.

4. Enjoy.
