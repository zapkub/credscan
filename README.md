# Credscan
minimal repository scanner service

## How Worker works
- find queued job if exist
- Clone repository to working dir
- iterate through the files
- write the result

## Todo
- [ ] Containerized
- [ ] handling repository scan error
- [ ] recover failed/crash job
- [ ] Optimize worker pool
