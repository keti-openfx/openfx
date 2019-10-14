kubectl delete -f ./yaml/job.yaml                                 
kubectl delete -f ./yaml/grafana-configmap.yaml                   
kubectl delete -f ./yaml/grafana.yml                              
kubectl delete -f ./yaml/prometheus.yml                           
sleep 30                              

kubectl apply -f ./yaml/prometheus.yml                            
sleep 30                                                          
                                                                  
kubectl apply -f ./yaml/grafana.yml                               
sleep 30                                                          
                                                                  
kubectl apply -f ./yaml/grafana-configmap.yaml                    
sleep 20 

kubectl apply -f ./yaml/job.yaml                                  
sleep 10                                                          
                                                                  
kubectl delete -f ./yaml/job.yaml                                 
sleep 10                                                          
                                                                  
kubectl apply -f ./yaml/job.yaml 
