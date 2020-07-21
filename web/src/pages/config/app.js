import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Row, Col, 
  Button, Popconfirm, Icon, message, Radio, Select, Switch,
  Transfer, Steps, Tag, Dropdown, Menu} from "antd";
import {connect} from "dva";
import Link from 'umi/link';
import {timeTrans, compareArray, hasPermission} from "@/utils/globalTools";
import Setup1 from './Ext1Setup1';
import Setup2 from './Ext1Setup2';
import Setup3 from './Ext1Setup3';
import AddSelect from './AddSelect';
import Ext1Form from './Ext1Form';
import Ext2Form from './Ext2Form';
import AppSync from './AppSync';
import styles from './index.module.css';
import {httpGet, httpPut} from '@/utils/request';

const FormItem = Form.Item;
const Option = Select.Option;
const RadioButton = Radio.Button;
const RadioGroup = Radio.Group;

@connect(({ loading, config, host, user }) => {
  return {
    projectsList: config.projectsList,
    projectsLoading: loading.effects['config/getProject'],
    projectsLen: config.projectsLen,
    configEnvList: config.configEnvList,
    appTypeList: config.appTypeList,
    deployTempleList:config.deployTempleList,
    robotList: user.robotList,
  }
})

class ProjectListPage extends React.Component {
  constructor(props) {
    super(props)
  }

  state = {
    loading: false,
    visible: false,
    addVisible: false,
    ext1Visible:false,
    ext2Visible: false,
    showTemplate: false,
    editCacheData: {},
    body: '',
    id: 0,
    info:{},
    loadDeployTempleList: 0,
    secondTable: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'config/getProject',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });
    dispatch({
      type: 'config/getAppType',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });
    dispatch({ 
      type: 'user/getRobot',
      payload: {
        page: 1,
        pageSize: 999, 
        status: 1,
      }
    });
  };


  getDeployExtend = (values) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/getDeployExtend',
      payload: {
        id: values
      }
    }).then(()=> {
      var tmp    = this.state.secondTable
      tmp[values]= this.props.deployTempleList
      this.setState({
        secondTable: tmp,
        loadDeployTempleList: values,
      })
    });
  }

  showEnvAddModal = () => {
    this.setState({
      isOk : false,
      editCacheData: {},
      visible: true,
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      info: {},
      deployTempleList: 0,
    });
  };

  cancelAddVisible = () => {
    this.setState({
      addVisible: false,
    });
  };

  handExt1Visible = () => {
    this.setState({
      addVisible:false,
      ext1Visible: true
    })
  }

  cancelExt1Visible = () => {
    this.setState({
      ext1Visible: false,
      info: {},
    })
  }

  cancelExt2Visible = () => {
    this.setState({
      ext2Visible: false,
      info: {},
    })
  }

  handExt2Visible = () => {
    this.setState({
      addVisible:false,
      ext2Visible: true
    })
  }

  // 关联角色是为了后期 deploy 权限检查
  // 数组比较
  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (Object.keys(obj).length) {
          values.Active      =  values.Active ? 1 : 0
          values.EnableSync  =  values.EnableSync ? 1 : 0
          values.DeployType  =  values.DeployType ? 1 : 0
          if (
            obj.Name          === values.Name &&
            obj.DeployType    === values.DeployType &&
            obj.EnableSync    === values.EnableSync &&
            obj.Active        === values.Active &&
            obj.Desc          === values.Desc &&
            obj.Tid           === values.Tid &&
            obj.EnvId         === values.EnvId &&
            obj.RepoUrl       === values.RepoUrl
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'config/configProjectEdit',
              payload: values,
            }).then(()=>{
              this.setState({ visible: false });
            });
          }
        } else {
          values.Active      = values.Active ? 1 : 0
          values.EnableSync  = values.EnableSync ? 1 : 0
          values.DeployType  =  values.DeployType ? 1 : 0
          dispatch({
            type: 'config/configProjectAdd',
            payload: values,
          }).then(()=>{
            this.setState({ visible: false })
          });
        }
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'config/configProjectDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };

  //显示编辑界面
  handleEdit = (values, page) => {
    values.title =  '编辑项目-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
  };

  // 显示项目发布配置模板，用于发布申请
  handleAppDeployConfig = (values) => {
    this.setState({ 
      ext1FromVisible: !this.state.ext1FromVisible,
      id: values.id,
      editCacheData: values,
    });
  };

  showAddSelect = (id) => {
    var tmp = this.state.info
    tmp["Aid"] = id
    this.setState({
      addVisible: true,
      info: tmp,
    })
  }

  showStepsAddModal = (id, info, isClone) => {
    if (info === undefined) {
      var tmp = this.state.info
      tmp["Aid"] = id
    } else {
      var tmp = info
      tmp["Aid"] = id
    }

    isClone && delete tmp.Dtid
    if (info.Extend == 1) {
      this.setState({
        ext1Visible: true,
        info: tmp,
      })
    } else{
      this.setState({
        ext2Visible: true,
        info: tmp,
      })
    }
  };

  // handleSync = (values) => {
  //   this.setState({ 
  //     showTemplate: true ,
  //     id: values,
  //   });
  // };

  onChange = (e) => {
  };

  onCheckChange = (checked) => {
  };

  switchTemplate = () => {
    this.setState({ 
      showTemplate: !this.state.showTemplate,
      id: 0,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'task/getProject',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  nextPage = () => {
    this.setState({ 
      
    });
  };

  handleDeployDelete = (id, dtid) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/deployExtendDel',
      payload: {
        id: id,
        tid: dtid,
      }
    });
  }

  handleSync = (record) => {
    this.setState({loading: true});
    httpGet(`/admin/sync/request/${record.id}`).then( res => {
      if (res.code != 200) {
        message.error(res.message)
        return
      }

      if (res.code == 200) {
         Modal.confirm({
          title: '确认初始项目？',
          content: `确认初始项目？`,
          onOk: () => {this.handleSyncConfirm(record.id)}
        })
      }
    }).finally(() => this.setState({loading: false}))
  }

  handleSyncConfirm = (values) => {
    this.setState({ 
      showTemplate: true ,
      id: values,
    });
  }

  columns = [
    {
      title: '项目名',
      dataIndex: 'Name',
    },{
      title: '发布环境',
      dataIndex: 'EnvId',
      'render': EnvId => this.props.configEnvList.map(x => {
        if (EnvId == x.id) {
          return x.Name
        }
      })
    },{
      title: '项目类型',
      dataIndex: 'Tid',
      'render': Tid => this.props.appTypeList.map(x => {
        if (Tid == x.id) {
          return x.Name
        }
      })
    },{
      title: '发布类型',
      dataIndex: 'DeployType',
      'render': DeployType => DeployType && '自定义发布' || '通用发布',
    },{
      title: '状态',
      dataIndex: 'Active',
      'render': Active => Active && '启用' || '禁用',
    }, {
      title: '备注',
      dataIndex: 'Desc',
      ellipsis: true,
    }, {
      title: '操作',
      key: 'action',
      render: (text, record) => (
        <span>
          {
            hasPermission('config-app-edit') && 
            <a onClick={()=>{this.handleEdit(record)}}>
              <Icon type="edit"/>编辑
            </a>
          }
          <Divider type="vertical" />
          {  
            hasPermission('config-app-set') && 
            // <a onClick={()=>{this.showStepsAddModal(record.id)}}>
            <a onClick={()=>{this.showAddSelect(record.id)}}>
              <Icon type="save"/>模板
            </a>
          }
          <Divider type="vertical" />
          {
            hasPermission('config-app-del') && 
            <Popconfirm title="你确定要删除吗?"  
              onConfirm={()=>{this.deleteRecord(record.id)}} 
              onCancel={()=>{this.cancel()}}
            >
              <a title="删除" >
                 <Icon type="delete" />删除
              </a>
            </Popconfirm>
          }
          {
            record.EnableSync == 1 &&
            <Divider type="vertical" />
          }
          {  record.EnableSync == 1 &&
            <Dropdown overlay={() => this.moreMenus(record)} trigger={['click']}>
              <a>
                更多 <Icon type="down"/>
              </a>
            </Dropdown>
          }
        </span>
      ),
    },
  ];

  moreMenus = (record) => (
    <Menu>
      <Menu.Item>
        {  
          hasPermission('config-app-set') && 
          <Link to={`/config/value?aid=${record.id}`}>
            <Icon type="setting"/>初始化变量
          </Link> 
        }
      </Menu.Item>
      <Menu.Divider/>
      <Menu.Item>
        {  
          hasPermission('config-app-init') && 
          <a onClick={()=>{this.handleSync(record)}} title="初始化" >
            <Icon type="redo" />初始化
          </a>
        }
      </Menu.Item>
    </Menu>
  );
  
  render() {
    const {visible, addVisible, ext1Visible, ext2Visible, 
      showTemplate, editCacheData,
      body, info, loadDeployTempleList, secondTable } = this.state;

    const {configEnvList, appTypeList, projectsList, projectsLoading, 
      projectsLen, hostsList, projectTargetKeys, imageList,
      deployTempleList,
      form: { getFieldDecorator} } = this.props;

    const addvar = <Button type="" onClick={this.showEnvAddModal} >添加应用</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('config-app-add') && <Col span={10}>{addvar}</Col>}
      </Row>;

    const expandedRowRender = (record) => {
      if (!secondTable[record.id]) {
        this.getDeployExtend(record.id)
      }

      // if (loadDeployTempleList!=record.id) {
      //   this.getDeployExtend(record.id)
      // }

      const columns = [
      {
        title: '模式',
        dataIndex: 'Extend',
        render: value => value == 1 ? <Icon style={{fontSize: 20, color: '#1890ff'}} type="ordered-list"/> :
          <Icon style={{fontSize: 20, color: '#1890ff'}} type="build"/>,
        width: 80
      }, 
      {
        title: '模板名字',
        dataIndex: 'TemplateName',
      }, {
        title: '发布审核',
        dataIndex: 'EnableCheck',
        render: value => value ? <Tag color="green">开启</Tag> : <Tag color="red">关闭</Tag>
      }, {
        title: '操作',
        render: info => (
          <span>
            { 
              hasPermission('config-app-edit') &&  
              <a onClick={()=>{this.showStepsAddModal(record.id, info)}}>编辑</a>
            }
            { 
              hasPermission('config-app-edit') && 
              <Divider type="vertical"/>
            }
            { 
              hasPermission('config-app-edit') && 
               <a onClick={()=>{this.showStepsAddModal(record.id, info, true)}}>克隆配置</a>
            }
            { 
              hasPermission('config-app-edit') && 
              <Divider type="vertical"/>
            }
            {
              hasPermission('config-app-del') && 
              <Popconfirm title="你确定要删除吗?"  
                onConfirm={()=>{this.handleDeployDelete(record.id, info.Dtid)}} 
                onCancel={()=>{this.cancel()}}
              >
                <a title="删除" >
                  删除
                </a>
              </Popconfirm>
            }
          </span>
        )
      }];

      return <Table
        rowKey="Dtid"
        loading={!secondTable[record.id]}
        columns={columns}
        dataSource={secondTable[record.id]}
        pagination={false} />
    };

    return (
      <div>
        <Modal
          title= { editCacheData.title || "新建项目" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
          style={{ top: 20 }} >
          <Form>
            <FormItem label="项目名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="项目类型">
              {getFieldDecorator('Tid', {
                initialValue: editCacheData.Tid || 'Please select' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  style={{ width: '100%' }}
                >
                  {appTypeList.map(x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem label="项目所属环境">
              <Col span={16}>
                {getFieldDecorator('EnvId', {
                  initialValue: editCacheData.EnvId || 'Please select' ,
                  rules: [{ required: true }],
                })(
                  <Select
                    placeholder="Please select"
                    style={{ width: '100%' }}
                  >
                    {configEnvList.map(x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
                  </Select>
                )}
              </Col>
              <Col span={6} offset={2}>
                <Link to="/config/environment">新建环境</Link>
              </Col>
            </FormItem>
            <FormItem label="是否需要初始化">
              {getFieldDecorator('EnableSync', {
                initialValue: editCacheData.EnableSync && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.EnableSync || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="是否生效">
              {getFieldDecorator('Active', {
                initialValue: editCacheData.Active && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.Active || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="是否需要自定义发布">
              {getFieldDecorator('DeployType', {
                initialValue: editCacheData.DeployType && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.DeployType || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input.TextArea />
              )}
            </FormItem>
          </Form> 
        </Modal>

        { showTemplate && 
          <AppSync
            id={this.state.id} 
            onCancel={this.switchTemplate} 
            onOk={body => this.setState({body})}
          />
        }

        { addVisible && 
          <AddSelect 
            cancelAddVisible={this.cancelAddVisible}
            ext1Visible={this.handExt1Visible}
            ext2Visible={this.handExt2Visible}
          />
        }

        { ext1Visible &&  
          <Ext1Form 
            cancelExt1Visible={this.cancelExt1Visible}
            editCacheData={this.state.editCacheData}
            configEnvList={configEnvList}
            info={this.state.info}
            dispatch={this.props.dispatch}
            robotList={this.props.robotList}
          />
        }

        { ext2Visible &&  
          <Ext2Form 
            cancelExt2Visible={this.cancelExt2Visible}
            editCacheData={this.state.editCacheData}
            configEnvList={configEnvList}
            info={this.state.info}
            dispatch={this.props.dispatch}
            robotList={this.props.robotList}
          />
        }

        <Card title="" extra={extra}>
          <Table 
            pagination={{
              showQuickJumper: true,
              total: projectsLen,
              showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
              onChange: this.pageChange
            }}
            columns={this.columns}
            dataSource={projectsList}
            expandedRowRender={expandedRowRender}
            loading={projectsLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(ProjectListPage);


