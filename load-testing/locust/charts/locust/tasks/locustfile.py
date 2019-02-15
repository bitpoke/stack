from locust import HttpLocust, TaskSet, task

class LocustTasks(TaskSet):
  @task
  def index(self):
      self.client.get("/")

class Locust(HttpLocust):
  task_set = LocustTasks
  min_wait = 1000
  max_wait = 3000
